# Anonimax

**Anonimax** é uma ferramenta de privacidade e anonimato para Linux que redireciona **todo o tráfego de rede do seu computador** pelo [Tor](https://www.torproject.org/), escondendo seu IP real e dificultando o rastreamento da sua navegação.

> Fork do [AnonSurf](https://gitlab.com/parrotsec/packages/anonsurf/) do ParrotSec, com funcionalidades extras de privacidade.

---

## Funciona no meu sistema?

| Sistema operacional | Funciona? |
|---|---|
| **Parrot OS** | ✅ Sim — sistema ideal, já tem quase tudo instalado |
| **Kali Linux** | ✅ Sim |
| **Ubuntu / Debian** | ✅ Sim, com algumas dependências extras |
| **Arch / Fedora / outros Linux** | ⚠️ Pode funcionar, mas sem suporte garantido |
| **Windows 10 / 11 / 7** | ❌ Não — a ferramenta usa tecnologias exclusivas do Linux |
| **macOS** | ❌ Não |

### Estou no Windows e quero privacidade — o que faço?

Se você usa Windows e quer algo parecido com o Anonimax, as melhores opções são:

- **[Tor Browser](https://www.torproject.org/download/)** — protege só o navegador, mas é fácil de instalar no Windows
- **[Tails OS](https://tails.boum.org)** — Linux completo que você bota em um pendrive USB e usa em qualquer PC sem instalar nada. É o que jornalistas e ativistas usam para anonimato real. Não deixa nenhum rastro no computador
- **[Whonix](https://www.whonix.org)** — máquina virtual Linux pré-configurada com Tor que roda dentro do Windows
- **[ProtonVPN](https://protonvpn.com)** — VPN gratuita confiável (não é Tor, mas esconde seu IP do provedor)

---

## O que o Anonimax faz?

Quando você executa `anonimax start`, o Anonimax:

1. 🔀 **Troca o endereço MAC** da sua placa de rede por um aleatório
2. 🌐 **Renova seu IP local** via DHCP
3. 🕐 **Força o fuso horário para UTC** (evita identificação pelo fuso)
4. 🔒 **Conecta uma VPN** opcional antes do Tor
5. 🧅 **Redireciona todo o tráfego pelo Tor**
6. 🚫 **Bloqueia WebRTC** (impede que o navegador vaze seu IP real)
7. 📡 **Bloqueia telemetria** do Google, Microsoft e Amazon
8. 📷 **Desativa câmera e microfone** via módulos do kernel

Quando você executa `anonimax stop`, tudo é restaurado ao estado original automaticamente.

---

## Comandos

```bash
anonimax start              # Ativa proteção completa
anonimax stop               # Para e restaura tudo
anonimax status             # Mostra o estado atual
anonimax restart            # Reinicia (troca circuito Tor)
anonimax changeid           # Troca de identidade Tor manualmente
anonimax autochange 10      # Troca de identidade a cada 10 minutos
anonimax autochange stop    # Para a troca automática
anonimax panic              # EMERGÊNCIA: derruba rede, limpa rastros e desliga
anonimax proxy set socks5://127.0.0.1:1080   # Configura proxy
anonimax proxy unset        # Remove proxy
anonimax config list        # Lista todas as configurações
anonimax config set block_telemetry true     # Altera uma configuração
```

### Modo Pânico

```bash
anonimax panic
```

Sequência de emergência: derruba todas as interfaces de rede → para o Tor → limpa as regras de firewall → restaura o MAC original → apaga rastros → **desliga o computador**. Use quando precisar fechar tudo rapidamente.

---

## Interface gráfica

O Anonimax tem uma interface gráfica (GUI) com todos os controles:

- Botão Start / Stop / Restart
- Painel de status (IP público, estado do Tor, DNS)
- Checkboxes para ativar/desativar cada funcionalidade individualmente

---

## O que é Tor?

O [Tor](https://www.torproject.org/) é uma rede gratuita e aberta que anonimiza o tráfego de internet passando-o por três servidores intermediários ao redor do mundo antes de chegar ao destino. Cada servidor conhece apenas o anterior e o próximo — nenhum deles vê quem você é e para onde está indo ao mesmo tempo.

O Tor Browser usa essa rede só para o navegador. O Anonimax vai além: coloca **todo o tráfego do sistema** pelo Tor.

---

## Funcionalidades detalhadas

| Funcionalidade | Comando / Config | Padrão |
|---|---|---|
| Aleatorizar MAC | `change_mac` | ✅ ativo |
| Renovar IP local | `renew_ip` | ✅ ativo |
| Forçar UTC | `force_utc` | ✅ ativo |
| Bloquear WebRTC | `block_webrtc` | ✅ ativo |
| Bloquear telemetria | `block_telemetry` | ❌ desativado |
| Desativar câmera/mic | `disable_camera_mic` | ❌ desativado |
| Limpeza ao parar | `clean_on_stop` | ❌ desativado |
| VPN antes do Tor | `use_vpn` | ❌ desativado |
| Proxy SOCKS5/HTTP | `use_proxy` | ❌ desativado |

Para ativar qualquer opção:
```bash
sudo anonimax config set block_telemetry true
```

Ou use os checkboxes na interface gráfica.

---

## Instalação

### Dependências

```bash
# Dependências de execução
sudo apt install tor iptables macchanger bleachbit network-manager

# Para VPN (opcional)
sudo apt install openvpn wireguard
```

### Opção 1 — Repositório APT (recomendado)

Configure uma vez e use `apt install` normalmente, igual a qualquer outro programa. Recebe atualizações automáticas com `apt upgrade`.

```bash
# 1. Adicionar a chave GPG do repositório
curl -fsSL https://revdf.github.io/anonimax/KEY.gpg \
  | sudo gpg --dearmor -o /usr/share/keyrings/anonimax.gpg

# 2. Adicionar o repositório
echo "deb [arch=amd64 signed-by=/usr/share/keyrings/anonimax.gpg] \
https://revdf.github.io/anonimax stable main" \
  | sudo tee /etc/apt/sources.list.d/anonimax.list

# 3. Instalar
sudo apt update
sudo apt install anonimax-cli          # CLI + daemon
sudo apt install anonimax-gui          # Interface gráfica (opcional)
```

Para remover o repositório depois:
```bash
sudo rm /etc/apt/sources.list.d/anonimax.list /usr/share/keyrings/anonimax.gpg
sudo apt update
```

### Opção 2 — Baixar o .deb manualmente

Se preferir instalar sem configurar o repositório, baixe os pacotes diretamente em:
**https://github.com/revdf/anonimax/releases**

```bash
sudo dpkg -i anonimax_6.1.0+anonimax1_amd64.deb
sudo dpkg -i anonimax-cli_6.1.0+anonimax1_amd64.deb
sudo dpkg -i anonimax-gui_6.1.0+anonimax1_amd64.deb   # opcional
sudo apt-get install -f                                  # corrige dependências
```

### Opção 3 — Compilar do código fonte

```bash
# Dependências de compilação
sudo apt install golang-go cmake g++ qt6-base-dev qt6-declarative-dev qml6-module-qtquick-shapes

# Clonar e instalar
git clone https://github.com/revdf/anonimax.git
cd anonimax
make build
sudo make install
```

---

## Créditos

Este projeto é um fork do **AnonSurf**, criado e mantido pela equipe do **[ParrotSec](https://www.parrotsec.org/)**:

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

## Aviso legal

Esta ferramenta é para uso em **segurança pessoal, privacidade e fins educacionais**. O uso para atividades ilegais é de responsabilidade exclusiva do usuário. Anonimato não é garantia absoluta de segurança — use com bom senso.
