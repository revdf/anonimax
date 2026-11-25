#!/usr/bin/env bash
# Anonimax Doctor — verifica se o ambiente está pronto para compilar e usar o Anonimax

RED='\033[0;31m'
GRN='\033[0;32m'
YLW='\033[1;33m'
BLD='\033[1m'
RST='\033[0m'

ok()   { echo -e "  ${GRN}✓${RST} $1"; }
fail() { echo -e "  ${RED}✗${RST} $1  ${RED}← FALTA${RST}"; MISSING+=("$2"); }
warn() { echo -e "  ${YLW}!${RST} $1"; }

MISSING=()

has_bin()  { command -v "$1" &>/dev/null; }
has_pkg()  { dpkg -l "$1" 2>/dev/null | grep -q '^ii'; }
has_any()  { has_bin "$1" || has_pkg "$1"; }

echo ""
echo -e "${BLD}Anonimax Doctor${RST}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# ─── Dependências de execução (CLI) ──────────────────────────────────────────
echo ""
echo -e "${BLD}[1] Dependências de execução (necessárias para usar o Anonimax)${RST}"

has_any tor        && ok "tor"              || fail "tor"              "tor"
has_any iptables   && ok "iptables"         || fail "iptables"         "iptables"
has_bin macchanger && ok "macchanger"       || fail "macchanger"       "macchanger"
has_any procps     && ok "procps"           || fail "procps"           "procps"
has_any systemd    && ok "systemd"          || fail "systemd"          "systemd"
has_any sudo       && ok "sudo"             || fail "sudo"             "sudo"

if has_bin dhclient || has_pkg isc-dhcp-client; then
    ok "dhclient (isc-dhcp-client)"
else
    fail "dhclient — renovar IP não vai funcionar" "isc-dhcp-client"
fi

# ─── Dependências opcionais ───────────────────────────────────────────────────
echo ""
echo -e "${BLD}[2] Dependências opcionais (features extras)${RST}"

has_bin bleachbit  && ok "bleachbit  (limpeza de rastros)"  || warn "bleachbit não instalado  — clean_on_stop e panic não vão limpar arquivos  →  sudo apt install bleachbit"
has_bin nyx        && ok "nyx        (monitor Tor)"         || warn "nyx não instalado  →  sudo apt install nyx"
has_bin openvpn    && ok "openvpn    (VPN)"                 || warn "openvpn não instalado  — use_vpn com OpenVPN não vai funcionar  →  sudo apt install openvpn"
has_bin wg         && ok "wireguard  (VPN)"                 || warn "wireguard não instalado  — use_vpn com WireGuard não vai funcionar  →  sudo apt install wireguard"

# ─── Dependências de compilação ───────────────────────────────────────────────
echo ""
echo -e "${BLD}[3] Dependências de compilação (necessárias para fazer o build)${RST}"

has_bin go         && ok "golang-go  ($(go version 2>/dev/null | awk '{print $3}'))"  || fail "golang-go"  "golang-go"
has_bin g++        && ok "g++        ($(g++ --version 2>/dev/null | head -1 | awk '{print $NF}'))"  || fail "g++"  "g++"
has_bin cmake      && ok "cmake      ($(cmake --version 2>/dev/null | head -1 | awk '{print $NF}'))"  || fail "cmake"  "cmake"
has_pkg qt6-base-dev           && ok "qt6-base-dev"           || fail "qt6-base-dev"           "qt6-base-dev"
has_pkg qt6-declarative-dev    && ok "qt6-declarative-dev"    || fail "qt6-declarative-dev"    "qt6-declarative-dev"
has_pkg qml6-module-qtquick-shapes && ok "qml6-module-qtquick-shapes" || fail "qml6-module-qtquick-shapes" "qml6-module-qtquick-shapes"

# ─── Resultado ───────────────────────────────────────────────────────────────
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if [ ${#MISSING[@]} -eq 0 ]; then
    echo -e "${GRN}${BLD}✓ Tudo pronto! Você pode compilar e usar o Anonimax.${RST}"
    echo ""
    echo "  Para compilar e instalar:"
    echo "    make build && sudo make install"
else
    echo -e "${YLW}${BLD}! Itens em falta: ${#MISSING[@]}${RST}"
    echo ""
    echo "  Instale tudo de uma vez com:"
    echo -e "    ${BLD}sudo apt install ${MISSING[*]}${RST}"
    echo ""
    echo "  Depois rode este script novamente para confirmar."
fi

echo ""
