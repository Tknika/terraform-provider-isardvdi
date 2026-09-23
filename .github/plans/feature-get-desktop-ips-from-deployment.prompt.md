# Plan: Exponer IP wireguard de los desktops de un deployment

## Contexto / hallazgos de investigación
- Provider usa terraform-plugin-framework. Resource clave: internal/provider/resource_isard_deployment.go, cliente: internal/client/deployment.go y internal/client/desktop.go.
- `Desktop` struct (internal/client/desktop.go) solo mapea id/name/description/template_id/vcpus/memory. `GetDesktop` golpea `GET /api/v3/domain/info/{id}` pero NO parsea ip ni viewer.
- En el código fuente real de IsardVDI (github.com/isard-vdi/isard), el documento de un domain incluye `viewer.guest_ip` (IP de wireguard, solo presente cuando el guest la reporta) y también un campo plano `ip` en varias proyecciones (`DeploymentDesktopsProcessed`, `UserDeploymentDesktop.ip`: "IP of the desktop. Only set when the desktop is started and has an IP assigned.").
- Existe endpoint apiv4 documentado: `GET /item/deployment/{deployment_id}/videowall` -> `DeploymentVideowallResponse.desktops: list[UserDeploymentDesktop]` con `id`, `name`, `status`, `ip`. En el webapp legacy (apiV3Segment) se llama igual pero montado bajo `/api/v3/item/deployment/{id}/videowall`.
- El `apiv3.json` (openapi/apiv3.json) que expone la instancia de prueba de Tknika (x.x.x.x) es un spec parcial/nuevo (FastAPI) que NO documenta `deployment` ni `domain/info` — esos son endpoints legacy no documentados ahí, pero SÍ son los que ya usa el provider en producción con éxito (`GetDeployment`, `GetDeploymentInfo`, `GetDesktopStatus` ya parsean `/api/v3/deployment/{id}`, `/api/v3/deployment/info/{id}`, `/api/v3/domain/info/{id}`).
- NO se pudo verificar remotamente el shape exacto de estos endpoints (IP x.x.x.x privada, no accesible desde este entorno; ?.tknika.eus/openapi devuelve 503). **Verificación pendiente manual por el usuario antes/durante la implementación** (ver Fase 0).
- `GetDeploymentInfo` ya existe y devuelve `map[string]interface{}` SIN parsear — es el candidato más barato a inspeccionar primero (puede que YA incluya desktops/ip y simplemente no se esté usando).
- No existe hoy ninguna función de cliente para: (a) listar los desktops de un deployment con su estado/IP, (b) esperar a que aparezca una IP.
- `StartDeployment` / `StopDeployment` / `WaitForDeploymentStopped` ya existen y seleccionan/arrancan TODOS los desktops del deployment — no hay (ni se necesita) arranque individual para este caso de uso.
- Decisiones ya acordadas con el usuario:
  1. El arranque + espera de IP debe ser **automático mediante un flag opcional** en `isardvdi_deployment` (ej. `wait_for_ip = true`), no una acción/recurso separado.
  2. El provider **solo expone las IPs** como atributo computed; NO orquesta Ansible (eso queda fuera, en ejemplos/documentación de referencia si acaso, no en el core).
  3. Verificación de endpoints: el usuario tiene una instancia de pruebas (x.x.x.x) pero no es accesible por mí; no hay swagger de apiv4 disponible. Toca verificar a mano.

## Fase 0 — Verificación de endpoints reales — ✅ RESUELTA, y REVERTIDA a v3 (decisión final, 23/09/2026)
**Cambio de rumbo importante tras revisión del usuario**: inicialmente se planteó usar apiv4 (`/api/v4/item/deployment/{id}/videowall` + `/start`) porque parecía el único sitio con el campo `ip`. Pero el usuario probó directamente el endpoint v3 que **ya usa el provider** (`GetDeployment` → `GET /api/v3/deployment/{id}`) contra la instancia de pruebas (IP *) con un deployment realmente arrancado, y la respuesta **ya incluye un array `desktops[]`** con exactamente lo que necesitamos:
```json
"desktops": [{
  "id": "58b2a320-...", "name": "Deployment Test", "state": "Started",
  "ip": "x.x.x.x", "user": "local-default-admin-admin", "user_name": "Administrator",
  ...
}]
```
Es decir: **no hace falta ningún endpoint nuevo ni apiv4**. El propio `GetDeployment` (ya usado en `Read()`, en `WaitForDeploymentStopped`, etc.) trae los desktops con su IP; simplemente el struct `DeploymentInfo` de hoy no mapea ese campo `desktops` y se descarta al deserializar.

**Decisión final de alcance**: implementar todo sobre **apiv3** (consistente con el resto del provider, sin añadir una segunda superficie de API que mantener). Se descarta por completo la vía apiv4 explorada antes (`GetDeploymentVideowall`, `StartDeploymentDesktopsV4` quedan sin implementar). Nota: el campo de estado en v3 se llama `state` (no `status` como en v4) — usar ese nombre al parsear.

- **Verificado en producción (?.tknika.eus) con datos reales, 23/09/2026** (esto sigue siendo válido como confirmación de que el concepto general —ip aparece cuando el desktop está arrancado, null/ausente si no— es correcto, aunque se hizo contra el endpoint v4 en su momento; el shape equivalente en v3 acaba de confirmarse igual de bien contra x.x.x.x).

## Fase 1 — Cliente (internal/client/deployment.go) — todo sobre v3, sin ficheros nuevos
3. Añadir a `DeploymentInfo` (mismo struct que ya rellena `GetDeployment`) un campo `Desktops []DeploymentDesktop \`json:"desktops"\`` y el nuevo tipo `DeploymentDesktop struct { ID, Name, State, IP, User, UserName string }` (json tags: `id`, `name`, `state`, `ip`, `user`, `user_name`; todos tolerantes a ausencia/null vía `omitempty` en el marshal, y el unmarshal de Go ya deja el zero-value string si el campo llega `null`/ausente). `GetDeployment` no cambia su lógica, solo el struct que rellena — el `json.Unmarshal` existente ya lo captura automáticamente.
4. `StartDeployment` (v3, ya existe, `PUT /api/v3/deployments/start/{id}`) se reutiliza tal cual — no hace falta ninguna función nueva de arranque.
5. Añadir `WaitForDeploymentDesktopsIP(deploymentID string, maxWaitSeconds int) (desktops []DeploymentDesktop, pendientes []string, err error)`: mismo patrón de polling que `WaitForDeploymentStopped` (ticker 2s, chequeo inmediato antes de esperar), pero llamando a `GetDeployment` (ya existente, sin tocar) y comprobando `len(d.IP) > 0` para todos los `deployment.Desktops`. Al agotar el timeout, en vez de devolver error, devuelve el último `desktops` obtenido + una lista `pendientes` con `name`/`id` de los que aún no tienen IP (para el warning, ver Further Considerations #2). Solo devuelve `error` si `GetDeployment` falla (p.ej. deployment no encontrado).
6. **No se toca nada de apiv4** — la carpeta `openapi/v4/` sigue sirviendo como referencia/documentación general de la API pero no se usa código contra ella en esta feature.

## Fase 2 — Resource `isardvdi_deployment` (internal/provider/resource_isard_deployment.go) — *depende de Fase 1*
6. Añadir al schema (`Schema()`):
   - `wait_for_ip` (`schema.BoolAttribute`, Optional+Computed, `booldefault.StaticBool(false)`): si es true, tras crear el deployment se llama a `StartDeployment` + `WaitForDeploymentDesktopsIP`.
   - `ip_wait_timeout_seconds` (`schema.Int64Attribute`, Optional+Computed, `int64default.StaticInt64(300)`).
   - `desktops` (`schema.ListNestedAttribute`, Computed): cada elemento con `id`, `name`, `ip`, `user`, `user_name` (todos `schema.StringAttribute{Computed: true}`; sin `status`).
7. Actualizar `deploymentResourceModel` con los campos `WaitForIP types.Bool`, `IPWaitTimeoutSeconds types.Int64`, `Desktops types.List` (tipo objeto `{id,name,ip,user,user_name}`).
8. `Create()` (tras la línea `plan.ID = types.StringValue(deploymentID)`): si `plan.WaitForIP.ValueBool()`, llamar `r.client.StartDeployment(deploymentID)` (la función v3 ya existente) y luego `r.client.WaitForDeploymentDesktopsIP(deploymentID, int(plan.IPWaitTimeoutSeconds.ValueInt64()))`; mapear el resultado (`id,name,ip,user,user_name`, descartando `state`) a `plan.Desktops` (usar `types.ListValueMust` con el `ObjectType` de desktop); si vuelve la lista de "sin IP a tiempo", emitir `resp.Diagnostics.AddWarning`. Si `wait_for_ip` es false, `plan.Desktops` queda como lista vacía (`types.ListValueMust(desktopObjType, []attr.Value{})`), nunca null, para evitar diffs de "unknown".
9. `Read()`: intentar refrescar `state.Desktops` llamando a `r.client.GetDeployment(id)` (que ya se llama aquí mismo, solo hay que leer también `deployment.Desktops` del resultado ya obtenido — no es una llamada HTTP adicional): si el deployment está parado (desktops sin IP), mantener lo que ya había en state en vez de vaciarlo — igual que se hace con vcpus/memory/network_interfaces.
10. `Update()`: si cambia `wait_for_ip` de false→true (o se solicita explícitamente), repetir el flujo de start+wait como en Create.

## Fase 3 — Documentación y ejemplos (*paralelo a Fase 2 una vez cerrado el diseño de atributos*)
11. Actualizar docs/resources/isardvdi_deployment.md: documentar `wait_for_ip`, `ip_wait_timeout_seconds`, `desktops` (lista de objetos `id`, `name`, `ip`, `user`, `user_name`) con nota clara de que las VMs quedan **encendidas** tras el apply si `wait_for_ip=true`, no se apagan automáticamente.
12. Añadir un ejemplo en examples/ (ej. `examples/deployment_with_wireguard_ip/main.tf`) mostrando `wait_for_ip = true` y un `output "desktops"` que expone `isardvdi_deployment.this.desktops` (con sus 5 campos), con un comentario indicando que ese output se puede consumir desde un inventory de Ansible (fuera del alcance del provider).

## Alcance excluido
- No se implementa arranque/parada automática tras el aprovisionamiento (usuario no lo pidió; las VMs quedan encendidas para permitir el acceso SSH/Ansible).
- No se integra Ansible en el provider (ni local-exec, ni provisioners) — solo se expone el dato.
- No se toca `isardvdi_vm` (desktop suelto) más allá de añadir `ip`/`status` al struct de cliente si aplica sin overhead; no es objetivo de este cambio añadir estos atributos a su schema salvo que el usuario lo pida luego.
- No se resuelve aún el endpoint exacto (Fase 0 es responsabilidad compartida: yo no tengo red hacia x.x.x.x ni acceso a ?.tknika.eus/openapi, que devolvió 503).

## Verificación
1. Fase 0 completada con curls reales y confirmación de shape JSON.
2. `go build ./...` y `go vet ./...` tras cada fase de código.
3. Prueba manual con `terraform apply` contra la instancia de pruebas: crear un `isardvdi_deployment` con `wait_for_ip = true`, comprobar que `terraform output`/state muestra IPs de wireguard reales y que son alcanzables por SSH.
4. Prueba de `terraform plan` posterior (sin cambios) para confirmar que no hay diffs espurios en `desktops` cuando `wait_for_ip = false`.

## Further Considerations
1. **Comparar v3 vs v4 a fondo**: descartado como parte de esta feature (ver Fase 0 — se resolvió todo con v3, que ya traía el dato). Sigue siendo válido como idea para el futuro si algún día se plantea modernizar el provider entero, pero no aporta nada a este cambio concreto.
2. **✅ CONFIRMADO por el usuario**: si el timeout se agota y solo algunas VMs tienen IP, NO se aborta el apply — se deja `ip = ""` para las que falten y se añade un `resp.Diagnostics.AddWarning(...)` listando qué desktops no consiguieron IP a tiempo.

## Campos a exponer — ✅ CONFIRMADO por el usuario (23/09/2026)
- **Nivel deployment**: `id`, `name`, `description` — ya existen como atributos top-level del resource (`deploymentResourceModel.ID/Name/Description`), no hace falta añadir nada nuevo a nivel deployment.
- **Nivel desktop** (dentro de la nueva lista `desktops`): `id`, `name`, `ip`, `user` (ID del propietario del desktop), `user_name` (nombre del propietario). **`status` NO se expone** como atributo público (el usuario no lo pidió), pero sigue siendo necesario internamente en `client.DeploymentDesktop` para la lógica de polling de `WaitForDeploymentDesktopsIP` (se puede dejar el campo Go sin mapear a un atributo del schema, o simplemente no incluirlo en el `ObjectType` de Terraform aunque exista en el struct Go).
- Esto ajusta la Fase 1 (struct `DeploymentDesktop` necesita también `User`, `UserName string`, parseados desde `desktops[].user` y `desktops[].user_name` en la respuesta de `GetDeploymentVideowall`) y la Fase 2 (el `ObjectType`/`ListNestedAttribute` de `desktops` en el schema pasa a `{id, name, ip, user, user_name}`, sin `status`).
