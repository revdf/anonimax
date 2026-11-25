# Anonimax

> Fork do [AnonSurf](https://gitlab.com/parrotsec/packages/anonsurf/) original do ParrotSec,
> com funcionalidades adicionais de privacidade e anonimato.

## Créditos ao projeto original

Este projeto é baseado no **AnonSurf**, criado e mantido pelo time do **ParrotSec**:

| Autor | Contribuição |
|---|---|
| Lorenzo "Palinuro" Faletra | Criador original do AnonSurf e do ParrotSec |
| Nong Hoang "DmKnght" Tu | Reescrita em Go, arquitetura atual |
| Lisetta "Sheireen" Ferrero | Contribuidora |
| Francesco "Mibofra" Bonanno | Contribuidor |
| Dario "danterolle" Camonita | Mantenedor atual |
| John "RightToPrivacy" Meagher | Documentação |
| Federica "marafed" Marasà | Ícone Nyx |

Repositório original: https://gitlab.com/parrotsec/packages/anonsurf/
Licença: GPL-3.0+

---

## O que este fork adiciona

Em relação ao AnonSurf original, este fork adiciona:

- **MAC Address aleatório** — troca na ativação, restaura ao parar
- **Renovação de IP local** — via dhclient no start/stop
- **VPN opcional** — conecta antes do Tor (OpenVPN ou WireGuard)
- **Bloquear WebRTC** — impede vazamento de IP real pelo navegador
- **Modo Pânico** (`anonsurf panic`) — derruba rede, limpa rastros, desliga
- **Limpeza automática ao parar** — cache, logs, arquivos temporários, DNS
- **Auto-troca de identidade** (`anonsurf autochange <N>`) — a cada N minutos
- **Forçar Timezone UTC** — evita identificação pelo fuso horário
- **Bloquear telemetria** — IPs conhecidos do Google, Microsoft etc.
- **Desativar microfone/câmera** — via módulos do kernel
- **Suporte a Proxy** — HTTP, HTTPS e SOCKS5, antes ou depois do Tor

Todas as opções podem ser configuradas tanto pelo **terminal** quanto pela **interface gráfica**.

---

## Comandos

```
anonsurf start            Inicia proteção completa
anonsurf stop             Para e restaura configurações originais
anonsurf status           Mostra estado de tudo
anonsurf changeid         Troca identidade Tor manualmente
anonsurf autochange <N>   Troca identidade a cada N minutos
anonsurf mac              Liga/desliga troca de MAC
anonsurf vpn              Configura VPN opcional
anonsurf ip               Renova IP local
anonsurf panic            Modo Pânico: para tudo e limpa rastros
anonsurf proxy set <addr> Configura proxy (ex: socks5://IP:PORTA)
anonsurf proxy unset      Remove proxy configurado
```

---

## Como compilar

Dependências de build: `cmake`, `g++`, `golang-go`, `git`, `qt6-base-dev`, `qt6-declarative-dev`, `qml6-module-qtquick-shapes`

Dependências de runtime: `bleachbit`, `nyx`, `sudo`, `tor`, `iptables`, `network-manager`, `procps`, `systemd`, `macchanger`, `dhclient`

```bash
git clone <este repositório>
cd anonsurf-enhanced && make build && sudo make install
make check
```

Para compilar apenas a GUI:
```bash
make gui
build/gui/anonsurf-gui
```

---

## Estrutura do código

- `src/cmd/` — pontos de entrada dos comandos Go
- `src/internal/anonsurf/` — orquestração do ciclo de vida
- `src/internal/dns/`, `firewall/`, `ipv6/`, `tor/` — domínios de rede
- `src/internal/command/`, `fileops/`, `systemd/` — adaptadores do SO
- `gui/` — interface Qt 6/QML
- `configs/` — arquivos de configuração

---

## Referências

- [Vivek Gite — dhclient hook script](https://www.cyberciti.biz/faq/dhclient-etcresolvconf-hooks/)
- [Bridges padrão do Tor](https://gitlab.torproject.org/tpo/anti-censorship/team/-/wikis/Default-Bridges)
- [Tails OS](https://tails.boum.org) — referência para Modo Pânico e limpeza de rastros
- [Whonix](https://www.whonix.org) — referência para isolamento de rede
