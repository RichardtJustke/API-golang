#!/bin/bash
# ============================================================
#  API TEST SUITE — http://localhost:8080
#  Uso: chmod +x api_tests.sh && ./api_tests.sh <modo>
#
#  Modos disponíveis:
#    attack        — flood de requisições GET (simulação de ataque)
#    concurrent    — criação simultânea de usuários
#    duplicados    — criação de usuários com dados repetidos
#    status        — mudança de estados em sequência
#    tudo          — roda attack + concurrent + status em sequência
#    limpar        — apaga todos os usuários com id > 20
# ============================================================

BASE="http://localhost:8080"
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3ODAxNTM4OTQsImlhdCI6MTc4MDE1MDI5NCwic3ViIjoiYXBpLWNsaWVudCJ9.RrGYKeecMWvTiwh0nvSuDcIZ1m8JsydhZlgadLNf368"

# ─── cores ───────────────────────────────────────────────────
RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'
CYAN='\033[0;36m'; BOLD='\033[1m'; NC='\033[0m'

log()  { echo -e "${CYAN}[INFO]${NC} $1"; }
ok()   { echo -e "${GREEN}[OK]${NC}   $1"; }
warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
err()  { echo -e "${RED}[ERRO]${NC} $1"; }

# ─── pega token ──────────────────────────────────────────────
get_token() {
  log "Obtendo token..."
  TOKEN=$(curl -s -X POST "$BASE/auth/token" | jq -r '.access_token // empty')
  if [[ -z "$TOKEN" ]]; then
    err "Falha ao obter token. A API está rodando?"
    exit 1
  fi
  ok "Token obtido: ${TOKEN:0:30}..."
}

auth_header() {
  echo "Authorization: Bearer $TOKEN"
}

# ============================================================
#  1. SIMULAÇÃO DE ATAQUE — flood de GET /users
# ============================================================
modo_attack() {
  local TOTAL=200
  local PARALELO=20
  log "=== ATAQUE: $TOTAL requisições GET /users ($PARALELO paralelas) ==="

  sucesso=0; falha=0
  inicio=$(date +%s%3N)

  attack_worker() {
    local res
    res=$(curl -s -o /dev/null -w "%{http_code}" \
      -H "$(auth_header)" \
      -H "Accept: application/json" \
      "$BASE/users?limit=10")
    echo "$res"
  }
  export -f attack_worker
  export TOKEN BASE

  for i in $(seq 1 $TOTAL); do
    attack_worker &
    # controla paralelismo
    if (( i % PARALELO == 0 )); then wait; fi
  done
  wait

  fim=$(date +%s%3N)
  duracao=$(( fim - inicio ))
  ok "Concluído em ${duracao}ms — verifique logs da API para erros 429/500."
}

# ============================================================
#  2. CRIAÇÃO SIMULTÂNEA — N usuários ao mesmo tempo
# ============================================================
modo_concurrent() {
  local QTDE=50
  log "=== CRIAÇÃO SIMULTÂNEA: $QTDE usuários em paralelo ==="

  criar_usuario() {
    local i=$1
    local ts=$(date +%s%N)
    local body="{
      \"nome\": \"Usuario$i\",
      \"cpf\": \"$(printf '%011d' $((RANDOM * RANDOM % 99999999999)))\",
      \"data_nascimento\": \"1990-01-01\",
      \"email\": \"user${i}_${ts}@test.com\",
      \"senha_hash\": \"hash_$i\",
      \"status_contrato\": \"ativo\",
      \"id_contrato\": \"$i\"
    }"
    local code
    code=$(curl -s -o /dev/null -w "%{http_code}" \
      -X POST "$BASE/users" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d "$body")
    echo "[$i] HTTP $code"
  }
  export -f criar_usuario
  export TOKEN BASE

  for i in $(seq 1 $QTDE); do
    criar_usuario "$i" &
  done
  wait
  ok "$QTDE requisições disparadas."
}

# ============================================================
#  3. USUÁRIOS DUPLICADOS — mesmo CPF e email várias vezes
# ============================================================
modo_duplicados() {
  local QTDE=20
  log "=== DUPLICADOS: $QTDE tentativas com mesmo CPF/email ==="
  log "Esperado: primeira = 201, demais = 409 Conflict"

  for i in $(seq 1 $QTDE); do
    code=$(curl -s -o /dev/null -w "%{http_code}" \
      -X POST "$BASE/users" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{
        "nome": "Duplicado Teste",
        "cpf": "00000000001",
        "data_nascimento": "1995-06-15",
        "email": "duplicado@test.com",
        "senha_hash": "hash_dup",
        "status_contrato": "ativo",
        "id_contrato": "DUP-001"
      }')
    if [[ "$code" == "201" ]]; then
      ok "[$i] Criado: HTTP $code"
    elif [[ "$code" == "409" ]]; then
      warn "[$i] Duplicado bloqueado: HTTP $code ✓"
    else
      err "[$i] Inesperado: HTTP $code"
    fi
  done
}

# ============================================================
#  4. MUDANÇA DE ESTADOS — ciclo em vários contratos
# ============================================================
modo_status() {
  local CONTRATOS=("1" "2" "3" "4" "5")
  local ESTADOS=("ativo" "suspenso" "cancelado" "ativo" "inutilizado" "ativo")
  log "=== MUDANÇA DE ESTADOS: ${#CONTRATOS[@]} contratos × ${#ESTADOS[@]} estados ==="

  for id in "${CONTRATOS[@]}"; do
    for estado in "${ESTADOS[@]}"; do
      resp=$(curl -s -w "\n%{http_code}" \
        -X PATCH "$BASE/users/contracts/status" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d "{\"id_contrato\":\"$id\",\"status\":\"$estado\"}")
      code=$(echo "$resp" | tail -1)
      body=$(echo "$resp" | head -1)
      if [[ "$code" == "200" ]]; then
        ok "Contrato $id → $estado: HTTP $code"
      else
        err "Contrato $id → $estado: HTTP $code | $body"
      fi
    done
  done
}

# ============================================================
#  5. TUDO — roda attack + concurrent + status
# ============================================================
modo_tudo() {
  log "=== MODO TUDO: executando todos os cenários ==="
  modo_attack
  echo ""
  modo_concurrent
  echo ""
  modo_status
  ok "=== TUDO CONCLUÍDO ==="
}

# ============================================================
#  6. LIMPAR — apaga usuários com id > 20
# ============================================================
modo_limpar() {
  log "=== LIMPEZA: apagando usuários com id > 20 ==="

  # busca todos os usuários
  todos=$(curl -s \
    -H "Authorization: Bearer $TOKEN" \
    -H "Accept: application/json" \
    "$BASE/users?limit=9999")

  ids=$(echo "$todos" | jq -r '.[] | .id' 2>/dev/null)

  if [[ -z "$ids" ]]; then
    warn "Nenhum usuário encontrado ou resposta inválida."
    echo "Resposta da API:"
    echo "$todos" | jq . 2>/dev/null || echo "$todos"
    exit 1
  fi

  apagados=0; mantidos=0
  for id in $ids; do
    if (( id > 20 )); then
      code=$(curl -s -o /dev/null -w "%{http_code}" \
        -X DELETE "$BASE/users/$id" \
        -H "Authorization: Bearer $TOKEN")
      if [[ "$code" == "200" || "$code" == "204" ]]; then
        ok "Deletado id=$id: HTTP $code"
        (( apagados++ ))
      else
        err "Falha ao deletar id=$id: HTTP $code"
      fi
    else
      log "Mantido id=$id"
      (( mantidos++ ))
    fi
  done

  echo ""
  ok "Resumo: $apagados apagados | $mantidos mantidos (ids 1–20)"
}

# ============================================================
#  MAIN
# ============================================================
MODE="${1:-help}"

if [[ "$MODE" == "help" || -z "$MODE" ]]; then
  echo -e "${BOLD}Uso:${NC} $0 <modo>"
  echo ""
  echo "  attack      Flood de GET /users (200 req, 20 paralelas)"
  echo "  concurrent  50 criações simultâneas de usuários"
  echo "  duplicados  20 tentativas com mesmo CPF/email"
  echo "  status      Ciclo de mudança de estados em contratos"
  echo "  tudo        Roda attack + concurrent + status"
  echo "  limpar      Apaga todos os usuários com id > 20"
  echo ""
  exit 0
fi

# verifica dependências
for dep in curl jq; do
  if ! command -v "$dep" &>/dev/null; then
    err "Dependência não encontrada: $dep  →  sudo apt install $dep"
    exit 1
  fi
done

get_token

case "$MODE" in
  attack)     modo_attack ;;
  concurrent) modo_concurrent ;;
  duplicados) modo_duplicados ;;
  status)     modo_status ;;
  tudo)       modo_tudo ;;
  limpar)     modo_limpar ;;
  *)          err "Modo inválido: $MODE"; exit 1 ;;
esac
