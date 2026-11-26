# Auditoria: Referências "AnonSurf" restantes no projeto Anonimax

## Contexto

Este projeto é o **Anonimax** — um fork do AnonSurf com identidade própria.
O objetivo desta tarefa é **apenas auditar e listar** todas as referências ao nome antigo que ainda restam. Não fazer nenhuma alteração.

---

## Tarefa

Varre **todas** as linhas de **todos** os arquivos do projeto (exceto as pastas abaixo) e lista cada ocorrência dos termos:

- `anonsurf` (qualquer capitalização: AnonSurf, ANONSURF, anon-surf, anon_surf)
- `Parrot Security team` (no contexto de autoria — deve virar `revdf`)
- `/usr/lib/anonsurf/` (caminho de binários)
- `/run/anonsurf/` (caminhos de estado em runtime)
- `/var/lib/anonsurf/` (caminhos de estado persistente)

**Excluir da busca:**
- `src/vendor/`
- `build/`
- `debian/`
- `.git/`
- Linhas com `package anonsurf` ou `import "…/anonsurf"` (nome do módulo Go interno)

---

## Como apresentar os resultados

Agrupe por categoria e apresente em tabela com colunas:
**Arquivo** | **Linha** | **Texto encontrado** | **Classificação**

### Classificações

| Classificação | Descrição |
|---------------|-----------|
| 🔴 `user-visible` | Aparece para o usuário: tooltip, mensagem de erro, notificação, label, nome de app |
| 🔴 `functional-path` | Caminho de arquivo ou serviço que afeta o funcionamento (`/usr/lib/anonsurf/`, `anonsurfd.service`) |
| 🟡 `internal-identifier` | Nome de classe C++, variável ou função interna — não aparece ao usuário |
| 🟡 `runtime-path` | Caminhos `/run/anonsurf/` ou `/var/lib/anonsurf/` — funcionam mas usam namespace antigo |
| ⚪ `comment-or-doc` | Comentário no código ou documentação — baixa prioridade |

---

## Extensões a verificar

`.go` `.cpp` `.h` `.qml` `.service` `.conf` `.desktop` `.install` `.menu` `.directory` `Makefile` `.cfg` `.md` `.txt`

---

## Resultado esperado

Lista completa e detalhada. **Não alterar nada** — apenas listar.
Ao final, apresentar contagem por categoria:
- Total 🔴 `user-visible`
- Total 🔴 `functional-path`
- Total 🟡 `internal-identifier`
- Total 🟡 `runtime-path`
- Total ⚪ `comment-or-doc`
