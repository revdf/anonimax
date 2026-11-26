# Tarefa: Substituir todas as referências "AnonSurf" por "Anonimax"

## Contexto

Este projeto é o **Anonimax** — um fork do AnonSurf com identidade própria.
Precisa-se substituir todas as referências ao nome antigo nos arquivos listados abaixo.
**Não alterar** arquivos fora desta lista e **não alterar** nomes de classes C++ internas (`AnonSurfBackend`, `AnonSurfConfig`, etc.) — apenas os itens explicitamente listados.

---

## Substituições a realizar

### 1. Strings visíveis ao usuário — QML

**`gui/qml/dashboard/ConnectionOverview.qml` linha 36**
```
"Managed by AnonSurf"  →  "Managed by Anonimax"
```

**`gui/qml/dashboard/ConnectionOverview.qml` linha 88**
```
"Starts AnonSurf during boot."  →  "Starts Anonimax during boot."
```

**`gui/qml/dashboard/ProtectionPanel.qml` linha 96**
```
"Removes AnonSurf routing and restores the previous network settings."
→
"Removes Anonimax routing and restores the previous network settings."
```

**`gui/qml/about/AboutWindow.qml` linha 62**
```
source: "qrc:/icons/anonsurf.png"  →  source: "qrc:/icons/anonimax.png"
```

---

### 2. Strings visíveis ao usuário — C++ (mensagens de erro/log)

**`gui/src/systemd/systemdjobs.cpp` linha 35**
```
"AnonSurf could not complete the %1 operation"
→
"Anonimax could not complete the %1 operation"
```

**`gui/src/cli/commandclient.cpp`**
```
"AnonSurf status is unavailable."           →  "Anonimax status is unavailable."
"AnonSurf returned an invalid status response."  →  "Anonimax returned an invalid status response."
"AnonSurf returned an incomplete status response."  →  "Anonimax returned an incomplete status response."
"AnonSurf returned an invalid Tor check response."  →  "Anonimax returned an invalid Tor check response."
```

**`gui/src/application/instance.cpp`**
```
"Cannot export the AnonSurf activation endpoint:"   →  "Cannot export the Anonimax activation endpoint:"
"Cannot register the AnonSurf desktop instance:"    →  "Cannot register the Anonimax desktop instance:"
"Cannot activate the existing AnonSurf window:"     →  "Cannot activate the existing Anonimax window:"
```

**`gui/src/application/main.cpp`**
```
app.setApplicationName("anonsurf-gui")  →  app.setApplicationName("anonimax-gui")
QIcon::fromTheme("anonsurf", QIcon(":/icons/anonsurf.png"))
→
QIcon::fromTheme("anonimax", QIcon(":/icons/anonimax.png"))
"The main AnonSurf window is unavailable."  →  "The main Anonimax window is unavailable."
```

**`gui/src/application/main.cpp` linha 96 (string de protectionDetail)**
```
"Waiting for AnonSurf status."  →  "Waiting for Anonimax status."
```
(esta string está em `gui/src/anonsurf/anonsurfbackend.h` linha 96)

---

### 3. Caminhos funcionais — systemd

**`gui/src/systemd/systemdclient.cpp`**
```
constexpr auto AnonSurfUnit = "anonsurfd.service";
→
constexpr auto AnonSurfUnit = "animaxd.service";

constexpr auto AnonSurfPath = "/org/freedesktop/systemd1/unit/anonsurfd_2eservice";
→
constexpr auto AnonSurfPath = "/org/freedesktop/systemd1/unit/animaxd_2eservice";
```

**`sys-units/anonsurf-identity.service`** — substituir no conteúdo do arquivo:
```
Requisite=anonsurfd.service   →  Requisite=animaxd.service
After=anonsurfd.service       →  After=animaxd.service
ExecStart=/usr/lib/anonsurf/anondaemon changeid  →  ExecStart=/usr/lib/anonimax/anondaemon changeid
```

**`sys-units/anonsurf-recovery.service`**
```
Before=anonsurfd.service      →  Before=animaxd.service
ExecStart=/usr/lib/anonsurf/anondaemon recover  →  ExecStart=/usr/lib/anonimax/anondaemon recover
```

**`sys-units/tor@default.service.d/anonsurf-recovery.conf`**
```
Wants=anonsurf-recovery.service  →  Wants=anonimax-recovery.service
After=anonsurf-recovery.service  →  After=anonimax-recovery.service
```

---

### 4. Caminho funcional — binário do daemon Go

**`src/internal/anonsurf/client.go` linha 10**
```
const BinaryPath = "/usr/lib/anonsurf/anondaemon"
→
const BinaryPath = "/usr/lib/anonimax/anondaemon"
```

---

### 5. Caminhos de runtime — estado em disco

Nos arquivos abaixo, substituir todos os caminhos `/run/anonsurf/` → `/run/anonimax/`
e `/var/lib/anonsurf/` → `/var/lib/anonimax/`:

- `src/internal/anonsurf/state.go`
- `src/internal/anonsurf/identity.go`
- `src/internal/proxy/proxy.go`
- `src/internal/firewall/manager.go`
- `src/internal/ipv6/manager.go`
- `src/internal/devices/camera_mic.go`
- `src/internal/dns/nm.go`
- `src/internal/dns/state.go`
- `src/internal/tor/torrc_manager.go`

Também em `src/internal/anonsurf/state.go` linha 19:
```
stateMarkerContent = "anonsurf\n"  →  stateMarkerContent = "anonimax\n"
```

E em `src/internal/dns/state.go` linha 17:
```
stateManaged = "anonsurf\n"  →  stateManaged = "anonimax\n"
```

---

### 6. Strings nos comandos CLI — Go (texto de help/exemplos)

**`src/cmd/anonsurf-cli/autochange.go`**
```
"  anonsurf autochange <N>     rotate every N minutes (N >= 1)\n"
→
"  anonimax autochange <N>     rotate every N minutes (N >= 1)\n"

"  anonsurf autochange stop    stop automatic rotation"
→
"  anonimax autochange stop    stop automatic rotation"

"Example: anonsurf autochange 10"  →  "Example: anonimax autochange 10"
```

**`src/cmd/anonsurf-cli/proxy.go`**
```
"Configure or remove the upstream proxy used with AnonSurf.\n"
→
"Configure or remove the upstream proxy used with Anonimax.\n"

"  anonsurf proxy set ..."  →  "  anonimax proxy set ..."   (todas as linhas)
"  anonsurf proxy unset"    →  "  anonimax proxy unset"
```

**`src/cmd/anonsurf-cli/status.go` linha 25**
```
"Run 'anonsurf status --nyx' for the interactive monitor or 'anonsurf status --json' for machine-readable output."
→
"Run 'anonimax status --nyx' for the interactive monitor or 'anonimax status --json' for machine-readable output."
```

**`src/cmd/anonsurf-cli/main.go` linha 25**
```
printer.Error(fmt.Sprintf("anonsurf: %v", err))
→
printer.Error(fmt.Sprintf("anonimax: %v", err))
```

**`src/cmd/dnstool/status.go`**
```
"read AnonSurf status: %w"                                   →  "read Anonimax status: %w"
"localhost is the only DNS server while AnonSurf is inactive" →  "localhost is the only DNS server while Anonimax is inactive"
```

---

### 7. README.md — atualizar exemplos de comandos

Substituir todos os exemplos de comando:
```
anonsurf start   →  anonimax start
anonsurf stop    →  anonimax stop
anonsurf status  →  anonimax status
(e assim por diante para todos os subcomandos)
sudo apt install anonsurf-cli  →  sudo apt install anonimax-cli
sudo apt install anonsurf-gui  →  sudo apt install anonimax-gui
sudo dpkg -i anonsurf_...      →  sudo dpkg -i anonimax_...
```

---

## O que NÃO alterar

- Nomes de classes C++: `AnonSurfBackend`, `AnonSurfConfig`, `AnonSurfUnit` (variável local), `AnonSurfBin`, `AnonSurfPath` (variáveis locais C++)
- Nomes de funções C++: `setAnonSurfEnabled`, `requestAnonSurfEnabled`, `anonsurfActive`
- Nomes de pacotes/imports Go: `package anonsurf`, `import "…/anonsurf"`
- Nomes de tipos Go: `AnonSurfActive` (campo de struct), `anonsurfActive` (variável local)
- O arquivo `sys-units/anonsurfd.service` (é o serviço original do fork — não instalado)
- A chain de iptables `ANONSURF-TELEMETRY` em `src/internal/firewall/telemetry.go` (nome técnico de regra de firewall)
- Referências ao AnonSurf original na seção de créditos do README (linha 183–187)

---

## Após fazer as substituições

Verificar se o projeto compila sem erros:
```bash
make build
```
