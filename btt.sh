#!/usr/bin/env bash

set -euo pipefail -o errtrace

notes="live2text"
default_host=http://127.0.0.1:8010/
#----

trap print_stack ERR
print_stack() {
    echo "❌ Error occurred at line $LINENO in script $0. Stack trace:"
    for i in "${!FUNCNAME[@]}"; do
        [[ $i -eq 0 ]] && continue
        echo "  at ${FUNCNAME[$i]}() in ${BASH_SOURCE[$i]}:${BASH_LINENO[$((i - 1))]}"
    done
}

# Helpers

merge_json() {
  jq -n 'reduce inputs as $item ({}; . *= $item)' < <(
    for json in "$@"; do
      echo -E "$json"
    done
  )
}

btt() {
  local args result
  args="$1"

  echo -E "$(date +"%Y-%m-%d %H:%M:%S.%N") Args: $args" >> ~/.btt/log
  result="$(osascript -e "tell application \"BetterTouchTool\" to ${args}")" || return $?
  echo -E "$(date +"%Y-%m-%d %H:%M:%S.%N") Res: $result" >> ~/.btt/log

  [[ "missing value" == "${result}" ]] && return 1
  [[ "success" == "${result}" || "deleted" == "${result}" ]] && return 0

  echo -E "${result}"
}

# Configs

json_icon() {
  local sf_symbol height only_icon json only_icon_json
  sf_symbol="$1"
  height="$2"
  only_icon="$3" # Format t/f — true/false

  json=$(cat <<EOF
  {
    "BTTTriggerConfig": {
      "BTTTouchBarItemSFSymbolDefaultIcon": "${sf_symbol}",
      "BTTTouchBarItemSFSymbolWeight": 0,
      "BTTTouchBarItemIconType": 2,
      "BTTTouchBarItemIconHeight": ${height},
      "BTTTouchBarItemPadding": -10
    }
  }
EOF
)
  only_icon_json=$(cat <<EOF
  {
    "BTTTriggerConfig": {
      "BTTTouchBarButtonColor": "0.0, 0.0, 0.0, 255.0",
      "BTTTouchBarOnlyShowIcon": true
    }
  }
EOF
)
  [[ "t" == "${only_icon}" ]] && json=$(merge_json "${json}" "${only_icon_json}")

  echo -E "${json}"
}

json_shell() {
  local fns payload interval action_type shell json json_trigger json_embed json_additional
  fns="$1"
  payload="$2"
  interval="$3"
  action_type="$4" # Format: n/e/a — none/embed/additional

  extract_header() {
    sed -En '1,/#----/p' <(btt "get_string_variable \"RECOGNIZER_SCRIPT\"")
  }
  extract_function() {
    sed -En "/^${1}\(\)/,/^}\$/p" <(btt "get_string_variable \"RECOGNIZER_SCRIPT\"")
  }

  shell="$(extract_header)
"
  while read -r fn; do
    shell+="$(extract_function "${fn}")
"
  done < <(echo "${fns}" | tr ' ' $'\n')
  shell+="${fns##* }"

  # hack to minimize the final script
  [[ "listen_socket" == "${fns}" ]] && shell="$(extract_function "${fns}" | sed '1d;$d')"

  shell="$(jq -Rs <<< "${shell}")"
  [[ -n "${payload}" ]] && shell="$(sed -e "${payload}" <<< "${shell}")"

  json_trigger=$(cat <<EOF
  {
    "BTTShellScriptWidgetGestureConfig": "\/bin\/bash:::-c:::-:::",
    "BTTTriggerConfig": {
      "BTTTouchBarAppleScriptStringRunOnInit": 1,
      "BTTTouchBarShellScriptString": ${shell}
    }
  }
EOF
)
  json_embed=$(cat <<EOF
  {
    "BTTPredefinedActionType": 206,
    "BTTShellTaskActionScript": ${shell},
    "BTTShellTaskActionConfig": "\/bin\/bash:::-c:::-:::"
  }
EOF
)
  json_additional=$(cat <<EOF
  {
    "BTTAdditionalActions": [${json_embed}]
  }
EOF
)

  [[ -n "${interval}" ]] && json_trigger=$(jq ".BTTTriggerConfig.BTTTouchBarScriptUpdateInterval = ${interval}" <<< "${json_trigger}")
  [[ "n" == "${action_type}" ]] && json="${json_trigger}"
  [[ "e" == "${action_type}" ]] && json="${json_embed}"
  [[ "a" == "${action_type}" ]] && json="${json_additional}"

  echo -E "${json}"
}

#json_separator() {
#  local json trigger
#  json=$(cat <<EOF
#  {
#    "BTTTriggerConfig": {
#      "BTTTouchBarButtonWidth": 8,
#      "BTTTouchBarButtonUseFixedWidth": 1
#    }
#  }
#EOF
#)
#  json=$(merge_json "$(json_trigger "Separator" "629 tb" 366 f)" "${json}")
#
#  echo -E "${json}"
#}

json_close() {
  local settings_group json
  settings_group="$1"

  json=$(cat <<EOF
  {
    "BTTOpenGroupWithName": "${settings_group}"
  }
EOF
)
  json=$(merge_json "$(json_trigger "Close Group" "629 tb" 205 f)" "$(json_icon "xmark.circle.fill" 25 t)" "${json}")

  echo -E "${json}"
}

json_trigger() {
  local name trigger action_type hidden json hidden_json
  name="$1"
  trigger="$2" # Format: type class tb/ot — BTTTriggerTypeTouchBar/BTTTriggerTypeOtherTriggers
  action_type="$3"
  hidden="$4" # Format: t/f — true/false

  json=$(cat <<EOF
  {
    "BTTTouchBarButtonName": "${name}",
    "BTTWidgetName": "${name}",
    "BTTTriggerName": "${name}",
    "BTTTriggerType": ${trigger%% *},
    "BTTTriggerClass": "$([[ "${trigger##* }" == "tb" ]] && echo "BTTTriggerTypeTouchBar" || echo "BTTTriggerTypeOtherTriggers")",
    "BTTPredefinedActionType": ${action_type},
    "BTTGroupName": "${notes}",
    "BTTNotes": "${notes}",
    "BTTTriggerConfig": {
      "BTTKeepGroupOpenWhileSwitchingApps": true
    }
  }
EOF
)
  [[ "${trigger##* }" == "ot" ]] && json="$(jq ".BTTGestureNotes = \"${notes}\"" <<< "${json}")"
  hidden_json=$(cat <<EOF
  {
    "BTTTriggerConfig": {
      "BTTTouchBarButtonWidth": 0,
      "BTTTouchBarButtonUseFixedWidth": 1
    }
  }
EOF
)
  [[ "t" == "${hidden}" ]] && json=$(merge_json "${json}" "${hidden_json}")

  echo -E "${json}"
}

# Templates

listen_socket() {
  local socket_path
  socket_path="<socket_path>"
  nc -U "${socket_path}" || echo "no subs"
}

listen_subs() {
  local host id
  id="<id>"

  host="${default_host}"
  subs="$(curl -sSLf "${host}api/subs" -H "Content-Type: application/json" -d "{\"id\": \"${id}\"}" || true)"
  [[ -z "${subs}" ]] && subs="no subs"
  echo "${subs}"
}

toggle_listen_socket() {
  local main_uuid default_interval is_start current_interval interval run_on_it json host
  main_uuid="<main_uuid>"
  default_interval="<default_interval>"
  is_start=0

  current_interval="$(btt "get_trigger \"${main_uuid}\"" | jq '.BTTTriggerConfig.BTTTouchBarScriptUpdateInterval')"
  if [[ 0 == "${current_interval}" ]]; then
    interval="${default_interval}"
    run_on_it=1
    is_start=1
  else
    interval=0
    run_on_it=0
    is_start=0
  fi


  host="${default_host}"
  local running_id
  running_id=$(btt "get_string_variable \"RECOGNIZER_ID\"" || true)
  if [[ -n "${running_id}" ]]; then
    curl -sSLf "${host}api/stop" -H "Content-Type: application/json" -d "{\"id\": \"${running_id}\"}" || true
    btt "set_persistent_string_variable \"RECOGNIZER_ID\" to \"\""
  fi

  json=$(cat <<EOF
  {
    "BTTTriggerConfig": {
      "BTTTouchBarAppleScriptStringRunOnInit": ${run_on_it},
      "BTTTouchBarScriptUpdateInterval": ${interval}
    }
  }
EOF
)

  # Stop it
  if [[ 0 == "${is_start}" ]]; then
    json="$(echo -E "${json}" | jq -c | jq -R)"
    btt "update_trigger \"${main_uuid}\" json $json" > /dev/null
    return 0
  fi

  # Run it
  local resp id socket_path
  device="$(btt "get_string_variable \"RECOGNIZER_DEVICE\"")"
  language="$(btt "get_string_variable \"RECOGNIZER_LANGUAGE\"")"
  resp=$(curl -sSLf "${host}api/start" -H "Content-Type: application/json" -d "{\"device\": \"${device}\", \"language\": \"${language}\"}" || true)
  [[ -z "${resp}" ]] && exit 1
  id=$(jq -r '.id' <<< "${resp}")
  socket_path=$(jq -r '.socketPath' <<< "${resp}")

  btt "set_persistent_string_variable \"RECOGNIZER_ID\" to \"${id}\""

  # curl implementation:
  # json=$(merge_json "${json}" "$(json_shell "listen_subs" "s/<main_uuid>/${main_uuid}/;s/<id>/${id}/" '0.25' 'n')")

  json=$(merge_json "${json}" "$(json_shell "listen_socket" "s|<socket_path>|${socket_path}|" '0.25' 'n')")
  json="$(echo -E "${json}" | jq -c | jq -R)"
  btt "update_trigger \"${main_uuid}\" json $json" > /dev/null
  return 0
}

select_device() {
  local device selected_device_uuid
  device="<device>"
  device="$(sed -E "s/ /%20/g" <<< "${device}")"

  selected_device_uuid=$(curl -sSLf "http://127.0.0.1:64444/get_string_variable/?variableName=RECOGNIZER_SELECTED_DEVICE_UUID")
  curl -sSLf "http://127.0.0.1:64444/set_string_variable/?variableName=RECOGNIZER_DEVICE&to=${device}"
  curl -sSLf "http://127.0.0.1:64444/refresh_widget/?uuid=${selected_device_uuid}"

#  selected_device_uuid="$(btt "get_string_variable \"RECOGNIZER_SELECTED_DEVICE_UUID\"")"
#  btt "set_persistent_string_variable \"RECOGNIZER_DEVICE\" to \"${device}\""
#  btt "refresh_widget \"${selected_device_uuid}\""
}

select_language() {
  local language selected_language_uuid
  language="<language>"

  selected_language_uuid=$(curl -sSLf "http://127.0.0.1:64444/get_string_variable/?variableName=RECOGNIZER_SELECTED_LANGUAGE_UUID")
  curl -sSLf "http://127.0.0.1:64444/set_string_variable/?variableName=RECOGNIZER_LANGUAGE&to=${language}"
  curl -sSLf "http://127.0.0.1:64444/refresh_widget/?uuid=${selected_language_uuid}"
}

open_settings() {
  local settings_group json
  settings_group="<settings_group>"

  json=$(cat <<EOF
  {
    "BTTPredefinedActionType": 205,
    "BTTOpenGroupWithName": "${settings_group}"
  }
EOF
)
  json="$(echo -E "${json}" | jq -c | jq -R)"
  btt "trigger_action $json" || true

  load_devices
}

# Actions

add_trigger() {
  local order parent_uuid json args
  order="$1"
  parent_uuid="$2"

  shift 2
  json="$(merge_json "$@" "{\"BTTOrder\":${order}}" | jq -c | jq -R)"
  args="add_new_trigger ${json}"
  [[ -n "${parent_uuid}" ]] && args="${args} parent_uuid \"${parent_uuid}\""

  btt "${args}"
}

delete_triggers() {
  local device_uuid uuids
  device_uuid="$(btt "get_string_variable \"RECOGNIZER_DEVICE_UUID\"")"
  uuids=$(btt "get_triggers trigger_parent_uuid \"${device_uuid}\" trigger_uuid 629" |
    jq -r '.[] | select((.BTTTriggerTypeDescription | IN("Close Group", "Separator")) | not) | .BTTUUID')

  for uuid in $uuids; do
    echo "uuid: $uuid"
#    btt "delete_triggers trigger_uuid \"${uuid}\""
  done
}

rerun() {
  stop

  cd ~/projects/live2text/
  go run ./... &
  pid=$!

  ok=$(osascript -e 'tell application "BetterTouchTool" to set_persistent_string_variable "RECOGNIZER_PID" to "'${pid}'"')
  [[ -z $(is_success "${ok}") ]] && echo "fail"
}

stop() {
  pid=$(osascript -e 'tell application "BetterTouchTool" to get_string_variable "RECOGNIZER_PID"')
  [[ -z $(has_value "${pid}") ]] && return
  kill -TERM "${pid}"
}

load_devices() {
  local host device_uuid exist_uuids devices uuid name

  host="$(btt "get_string_variable \"RECOGNIZER_HOST\"")"
  device_uuid="$(btt "get_string_variable \"RECOGNIZER_DEVICE_UUID\"")"
  exist_uuids=$(btt "get_triggers trigger_parent_uuid \"${device_uuid}\" trigger_id 629" |
    jq -r '.[] | select((.BTTTriggerTypeDescription | IN("Close Group", "Separator")) | not) | .BTTUUID + "\t" + .BTTTouchBarButtonName')

  devices=$(curl -sSL "${host}api/devices" || true)
  [[ -z "${devices}" ]] && exit 1
  devices=$(jq -r '.devices[]' <<< "${devices}")

  [[ -n "${exist_uuids}" ]] && while IFS=$'\t' read -r uuid name; do
    grep -Eq "^${name}\$" <<< "$devices" || btt "delete_trigger \"${uuid}\""
  done <<< "${exist_uuids}"

  [[ -n "${devices}" ]] && while IFS=$'\t' read -r name; do
    grep -Eq "\t${name}\$" <<< "${exist_uuids}" ||
      add_trigger 3 "${device_uuid}" \
      "$(json_trigger "${name}" "629 tb" 366 f)" \
      "$(json_shell "btt select_device" "s/<device>/${name}/" '' 'a')"
  done <<< "${devices}"

  return
}

print_selected_device() {
  local device
  device=$(btt "get_string_variable \"RECOGNIZER_DEVICE\"" || true)
  [[ -z "${device}" ]] && echo "⚠️ No device selected" || echo "✅ ${device}"
}

print_selected_language() {
  local language
  language=$(btt "get_string_variable \"RECOGNIZER_LANGUAGE\"" || true)
  [[ -z "${language}" ]] && echo "⚠️ No language selected" || echo "✅ ${language}"
}

print_status() {
  local host ok
  host=$(btt "get_string_variable \"RECOGNIZER_HOST\"")
  ok=$(curl -sSL "${host}api/health" 2>/dev/null || true)
  [[ -n "${ok}" ]] && echo "✅" || echo "⚠️"
}

get_trigger() {
  local uuid
  uuid="$1"

  btt "get_trigger \"${uuid}\""
}

# Shell context
init() {
  clear

  local subs_uuid device_uuid selected_device_uuid language_uuid selected_language_uuid script_path settings_group socket_path
  script_path="$(realpath "$0")"
  settings_group=Subs
  socket_path="\/tmp\/recognizer"

  # Set Hosts
  btt "set_persistent_string_variable \"RECOGNIZER_HOST\" to \"${default_host}\""
  btt "set_persistent_string_variable \"RECOGNIZER_SCRIPT\" to $(jq -Rs < "${script_path}")"

  # Subs Group
  subs_uuid="$(add_trigger 10 "" "$(json_trigger "${settings_group}" "630 tb" 206 t)")"
  btt "set_persistent_string_variable \"RECOGNIZER_SUBS_UUID\" to \"${subs_uuid}\""
  add_trigger 0 "${subs_uuid}" "$(json_trigger "Close Group" "629 tb" 191 f)" "$(json_icon "xmark.circle.fill" 25 t)"
  add_trigger 1 "${subs_uuid}" "$(json_trigger "⏳" "642 tb" 366 f)" "$(json_shell "btt print_status" '' 15 'n')"

  # Device Group
  device_uuid=$(add_trigger 2 "${subs_uuid}" "$(json_trigger "Device" "630 tb" 206 f)" "$(json_icon "microphone" 22 f)") #"$(json_shell "btt merge_json json_trigger add_trigger load_devices" '' '' 'e')")
  btt "set_persistent_string_variable \"RECOGNIZER_DEVICE_UUID\" to \"${device_uuid}\""
  add_trigger 0 "${device_uuid}" "$(json_close "${settings_group}")"
  selected_device_uuid="$(add_trigger 1 "${device_uuid}" \
    "$(json_trigger "⏳ Selected Device" "642 tb" 366 f)" \
    "$(json_shell "btt print_selected_device" '' 15 'n')" \
    "{\"BTTTriggerConfig\":{\"BTTTouchBarFreeSpaceAfterButton\":25}}"
  )"
  btt "set_persistent_string_variable \"RECOGNIZER_SELECTED_DEVICE_UUID\" to \"${selected_device_uuid}\""

  # Language Group
  language_uuid="$(add_trigger 3 "${subs_uuid}" "$(json_trigger "Language" "630 tb" 206 f)" "$(json_icon "character" 22 f)")"
  btt "set_persistent_string_variable \"RECOGNIZER_LANGUAGE_UUID\" to \"${language_uuid}\""
  add_trigger 0 "${language_uuid}" "$(json_close "${settings_group}")"
  selected_language_uuid="$(add_trigger 1 "${language_uuid}" \
    "$(json_trigger "⏳ Selected Language" "642 tb" 366 f)" \
    "$(json_shell "btt print_selected_language" '' 15 'n')" \
    "{\"BTTTriggerConfig\":{\"BTTTouchBarFreeSpaceAfterButton\":25}}"
  )"
  btt "set_persistent_string_variable \"RECOGNIZER_SELECTED_LANGUAGE_UUID\" to \"${selected_language_uuid}\""
  local i=3
  for language in en-US es-ES fr-FR pt-BR ru-RU ja-JP de-DE; do
    add_trigger $((i++)) "${language_uuid}" \
      "$(json_trigger "${language}" "629 tb" 366 f)" \
      "$(json_shell "btt select_language" "s/<language>/${language}/" '' 'a')"
  done

  # Named & main
  add_trigger 0 "" \
    "$(json_trigger "${settings_group}" "643 ot" 366 f)" \
    "$(json_shell "btt extract_header extract_function merge_json json_trigger json_shell add_trigger load_devices open_settings" "s/<settings_group>/${settings_group}/" '' 'a')"

  local main_uuid
  main_uuid="$(uuidgen)"
  add_trigger 11 "" \
    "$(json_trigger "Speech recognition" "642 tb" 366 f)" \
    "$(json_shell "btt merge_json json_shell toggle_listen_socket" "s/<main_uuid>/${main_uuid}/;s/<default_interval>/0.25/" '' 'e')" \
    "$(json_shell "listen_subs" "s/<main_uuid>/${main_uuid}/;s/<id>/SOME_ID/" '0' 'n')" \
    "{\"BTTTriggerConfig\":{\"BTTTouchBarLongPressActionName\": \"${settings_group}\"}}" \
    "{\"BTTTriggerConfig\":{\"BTTTouchBarButtonFontSize\":12,\"BTTTouchBarButtonColor\":\"0.0, 0.0, 0.0, 255.0\",\"BTTTouchBarButtonTextAlignment\":0}}" \
    "{\"BTTUUID\":\"${main_uuid}\"}"
}

clear() {
  local uuids
  btt 'trigger_action "{\"BTTPredefinedActionType\":191}"' || true

  uuids=$(btt "get_triggers" | jq -r ".[] | select(.BTTGroupName == \"${notes}\" or .BTTNotes == \"${notes}\" or .BTTGestureNotes == \"${notes}\") | .BTTUUID")
  while read -r uuid; do
    btt "delete_trigger \"${uuid}\""
  done <<< "${uuids}"
}

[[ $# -lt 1 ]] && echo "No function passed" && exit 1
fn="$1"
shift
"${fn}" "$@"