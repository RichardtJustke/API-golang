#!/bin/bash
# ============================================================
#  STRESS TEST EXTREMO — http://localhost:8080
#  Uso: chmod +x stress_test.sh && ./stress_test.sh <modo>
#
#  Modos:
#    attack      — 5000 GET /users (200 paralelas)
#    concurrent  — 200 criações simultâneas
#    duplicados  — 50 tentativas com mesmo CPF/email
#    status      — ciclo de estados em 20 contratos
#    mixed       — GET + POST + PATCH ao mesmo tempo
#    tudo        — roda todos em sequência
#    limpar      — apaga usuários com id > 20
# ============================================================

BASE="http://localhost:8080"
TOKEN=""

RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'
CYAN='\033[0;36m'; BOLD='\033[1m'; NC='\033[0m'

log()  { echo -e "${CYAN}[INFO]${NC} $1"; }
ok()   { echo -e "${GREEN}[OK]${NC}   $1"; }
warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
err()  { echo -e "${RED}[ERRO]${NC} $1"; }

get_token() {
  log "Obtendo token..."
  TOKEN=$(curl -s -X POST "$BASE/auth/token" | jq -r '.access_token // empty')
  if [[ -z "$TOKEN" ]]; then
    err "Falha ao obter token. A API está rodando?"
    exit 1
  fi
  ok "Token obtido: ${TOKEN:0:30}..."
}

# ============================================================
#  MÉTRICAS
# ============================================================
declare -A RESULTADOS
contar_resultado() { RESULTADOS[$1]=$(( ${RESULTADOS[$1]:-0} + 1 )); }

imprimir_resumo() {
  local total=0
  echo ""
  echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
  echo -e "${BOLD}  RESUMO DE RESULTADOS${NC}"
  echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
  for code in $(echo "${!RESULTADOS[@]}" | tr ' ' '\n' | sort); do
    count=${RESULTADOS[$code]}
    total=$(( total + count ))
    if   [[ "$code" == 2* ]]; then echo -e "  ${GREEN}HTTP $code${NC}: $count req"
    elif [[ "$code" == 4* ]]; then echo -e "  ${YELLOW}HTTP $code${NC}: $count req"
    elif [[ "$code" == 5* ]]; then echo -e "  ${RED}HTTP $code${NC}: $count req"
    else echo -e "  HTTP $code: $count req"; fi
  done
  echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
  echo -e "  Total: ${BOLD}$total${NC} requisições"
  echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
  RESULTADOS=()
}

# ============================================================
#  1. ATAQUE — 5000 GET /users, 200 paralelas
# ============================================================
modo_attack() {
  local TOTAL=5000 PARALELO=200
  log "=== ATAQUE EXTREMO: $TOTAL req GET /users ($PARALELO paralelas) ==="
  local inicio=$(date +%s%3N)
  local tmp=$(mktemp)

  worker_attack() {
    curl -s -o /dev/null -w "%{http_code}\n" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Accept: application/json" \
      "$BASE/users?limit=10"
  }
  export -f worker_attack
  export TOKEN BASE

  seq 1 $TOTAL | xargs -P $PARALELO -I{} bash -c 'worker_attack' >> "$tmp"

  local fim=$(date +%s%3N)
  local dur=$(( fim - inicio ))

  while IFS= read -r code; do contar_resultado "$code"; done < "$tmp"
  rm "$tmp"

  local rps=$(( TOTAL * 1000 / dur ))
  ok "Concluído em ${dur}ms (~${rps} req/s)"
  imprimir_resumo
}

# ============================================================
#  2. CRIAÇÃO SIMULTÂNEA — 200 usuários em paralelo
# ============================================================
modo_concurrent() {
  local QTDE=200
  log "=== CRIAÇÃO SIMULTÂNEA: $QTDE usuários em paralelo ==="
  local tmp=$(mktemp)

  criar_usuario() {
    local i=$1 ts=$(date +%s%N)
    local cpf=$(printf '%011d' $(( (RANDOM * RANDOM * RANDOM) % 99999999999 )))
    local body="{\"nome\":\"StressUser$i\",\"cpf\":\"$cpf\",\"data_nascimento\":\"1990-06-15\",\"email\":\"stress${i}_${ts}@test.com\",\"senha_hash\":\"hash$i\",\"status_contrato\":\"ativo\",\"id_contrato\":\"STRESS-$i\"}"
    curl -s -o /dev/null -w "%{http_code}\n" \
      -X POST "$BASE/users" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d "$body"
  }
  export -f criar_usuario
  export TOKEN BASE

  seq 1 $QTDE | xargs -P $QTDE -I{} bash -c 'criar_usuario "$@"' _ {} >> "$tmp"

  while IFS= read -r code; do contar_resultado "$code"; done < "$tmp"
  rm "$tmp"

  ok "$QTDE requisições disparadas."
  imprimir_resumo
}

# ============================================================
#  3. DUPLICADOS — 50 tentativas com mesmo CPF/email
# ============================================================
modo_duplicados() {
  local QTDE=50
  log "=== DUPLICADOS: $QTDE tentativas com mesmo CPF/email ==="
  local tmp=$(mktemp)

  dup_worker() {
    curl -s -o /dev/null -w "%{http_code}\n" \
      -X POST "$BASE/users" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{"nome":"DupStress","cpf":"00000000099","data_nascimento":"1995-01-01","email":"dupstress@test.com","senha_hash":"hash_dup","status_contrato":"ativo","id_contrato":"DUP-STRESS"}'
  }
  export -f dup_worker
  export TOKEN BASE

  seq 1 $QTDE | xargs -P $QTDE -I{} bash -c 'dup_worker' >> "$tmp"

  while IFS= read -r code; do contar_resultado "$code"; done < "$tmp"
  rm "$tmp"

  ok "$QTDE requisições disparadas."
  imprimir_resumo
}

# ============================================================
#  4. MUDANÇA DE ESTADOS — 20 contratos × 5 estados
# ============================================================
modo_status() {
  local CONTRATOS=(1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20)
  local ESTADOS=("ativo" "suspenso" "cancelado" "suspenso" "ativo")
  log "=== ESTADOS: ${#CONTRATOS[@]} contratos × ${#ESTADOS[@]} estados ==="
  local tmp=$(mktemp)

  status_worker() {
    local id=$1 estado=$2
    curl -s -o /dev/null -w "%{http_code}\n" \
      -X PATCH "$BASE/users/contracts/status" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d "{\"id_contrato\":\"$id\",\"status\":\"$estado\"}"
  }
  export -f status_worker
  export TOKEN BASE

  for id in "${CONTRATOS[@]}"; do
    for estado in "${ESTADOS[@]}"; do
      status_worker "$id" "$estado" >> "$tmp" &
    done
  done
  wait

  while IFS= read -r code; do contar_resultado "$code"; done < "$tmp"
  rm "$tmp"

  ok "Estados aplicados."
  imprimir_resumo
}

# ============================================================
#  5. MIXED — GET + POST + PATCH ao mesmo tempo
# ============================================================
modo_mixed() {
  local QTDE=500
  log "=== MIXED: $QTDE req misturadas (GET + POST + PATCH) em paralelo ==="
  local tmp=$(mktemp)

  mixed_worker() {
    local i=$1
    local mod=$(( i % 3 ))
    local ts=$(date +%s%N)
    if [[ $mod -eq 0 ]]; then
      curl -s -o /dev/null -w "%{http_code}\n" \
        -H "Authorization: Bearer $TOKEN" "$BASE/users?limit=5"
    elif [[ $mod -eq 1 ]]; then
      local cpf=$(printf '%011d' $(( (RANDOM * RANDOM) % 99999999999 )))
      curl -s -o /dev/null -w "%{http_code}\n" \
        -X POST "$BASE/users" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d "{\"nome\":\"Mix$i\",\"cpf\":\"$cpf\",\"data_nascimento\":\"1990-01-01\",\"email\":\"mix${i}_${ts}@test.com\",\"senha_hash\":\"h$i\",\"status_contrato\":\"ativo\",\"id_contrato\":\"MIX-$i\"}"
    else
      local estado
      case $(( i % 3 )) in
        0) estado="ativo" ;; 1) estado="suspenso" ;; *) estado="cancelado" ;;
      esac
      curl -s -o /dev/null -w "%{http_code}\n" \
        -X PATCH "$BASE/users/contracts/status" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d "{\"id_contrato\":\"$(( i % 20 + 1 ))\",\"status\":\"ativo\"}"
    fi
  }
  export -f mixed_worker
  export TOKEN BASE

  seq 1 $QTDE | xargs -P 100 -I{} bash -c 'mixed_worker "$@"' _ {} >> "$tmp"

  while IFS= read -r code; do contar_resultado "$code"; done < "$tmp"
  rm "$tmp"

  ok "$QTDE requisições mistas disparadas."
  imprimir_resumo
}

# ============================================================
#  6. TUDO
# ============================================================
modo_tudo() {
  log "=== MODO TUDO ==="
  modo_attack
  modo_concurrent
  modo_duplicados
  modo_status
  modo_mixed
  ok "=== STRESS TEST CONCLUÍDO ==="
}

# ============================================================
#  7. LIMPAR
# ============================================================
modo_limpar() {
  log "=== LIMPEZA: apagando usuários com id > 20 ==="
  todos=$(curl -s -H "Authorization: Bearer $TOKEN" -H "Accept: application/json" "$BASE/users?limit=9999")
  ids=$(echo "$todos" | jq -r '.[] | .id' 2>/dev/null)

  if [[ -z "$ids" ]]; then
    warn "Nenhum usuário encontrado."
    exit 1
  fi

  apagados=0; mantidos=0
  for id in $ids; do
    if (( id > 20 )); then
      code=$(curl -s -o /dev/null -w "%{http_code}" -X DELETE "$BASE/users/$id" -H "Authorization: Bearer $TOKEN")
      if [[ "$code" == "200" || "$code" == "204" ]]; then
        ok "Deletado id=$id"
        (( apagados++ ))
      else
        err "Falha id=$id: HTTP $code"
      fi
    else
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
  echo "  attack      5000 GET /users (200 paralelas)"
  echo "  concurrent  200 criações simultâneas"
  echo "  duplicados  50 tentativas com mesmo CPF/email"
  echo "  status      Ciclo de estados em 20 contratos"
  echo "  mixed       GET + POST + PATCH misturados (500 req)"
  echo "  tudo        Roda todos os modos"
  echo "  limpar      Apaga usuários com id > 20"
  echo ""
  exit 0
fi

for dep in curl jq xargs; do
  if ! command -v "$dep" &>/dev/null; then
    err "Dependência não encontrada: $dep"
    exit 1
  fi
done

get_token

case "$MODE" in
  attack)     modo_attack ;;
  concurrent) modo_concurrent ;;
  duplicados) modo_duplicados ;;
  status)     modo_status ;;
  mixed)      modo_mixed ;;
  tudo)       modo_tudo ;;
  limpar)     modo_limpar ;;
  *)          err "Modo inválido: $MODE"; exit 1 ;;
esac
