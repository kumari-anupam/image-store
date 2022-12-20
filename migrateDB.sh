#! /bin/bash

urlencode() {
    old_lc_collate=$LC_COLLATE
    LC_COLLATE=C

    local length="${#1}"
    for i in $(seq 0 $((length-1))); do
        local c="${1:$i:1}"
        case $c in
            [a-zA-Z0-9.~_-]) printf '%s' "$c" ;;
            *) printf '%%%02X' "'$c" ;;
        esac
    done

    LC_COLLATE=$old_lc_collate
}

urldecode() {
    local url_encoded="${1//+/ }"
    printf '%b' "${url_encoded//%/\\x}"
}

DB_USERNAME_ENC=`urlencode $DB_USERNAME`
DB_PASSWORD_ENC=`urlencode $DB_PASSWORD`

go get "github.com/migrate"
migrate -source "file://${DB_MIGRATION_PATH}" \
        -database "${DB_DRIVER}"://"${DB_USERNAME_ENC}":"${DB_PASSWORD_ENC}"@"${DB_HOST}:${DB_PORT}"/"${DB_NAME}"?sslmode=disable up