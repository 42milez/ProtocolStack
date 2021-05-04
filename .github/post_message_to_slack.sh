#!/bin/bash

set -eu

read -r COLOR STATUS_MESSAGE < <(
case "${NEEDS_PREVIOUS_JOB_RESULT}" in
  "success") echo "#74c7b8" "passed";;
  "failure") echo "#ef4f4f" "failed";;
  "cancelled") echo "#f4d160" "was cancelled";;
  "skipped") echo "#dddddd" "was skipped";;
  *) exit 1
esac
)
readonly COLOR
readonly STATUS_MESSAGE

readonly COMMIT_HASH=$(echo "${GITHUB_SHA}" | cut -c 1-7)

readonly LINK=$(
if [[ "${GITHUB_REPOSITORY_REF}" =~ ^"refs/heads/" ]]; then
  readonly BRANCH=$(echo "${GITHUB_REPOSITORY_REF}" | cut -c 12-)
  echo "branch: <https://github.com/${GITHUB_REPOSITORY}/tree/${BRANCH}|${BRANCH}>"
else
  readonly TAG=$(echo "${GITHUB_REPOSITORY_REF}" | cut -c 11-)
  echo "tag: <https://github.com/${GITHUB_REPOSITORY}/releases/tag/${TAG}|${TAG}>"
fi
)

readonly TEXT=$(cat <<EOF
Workflow: ${GITHUB_WORKFLOW} (<https://github.com/${GITHUB_REPOSITORY}/actions/runs/${GITHUB_RUN_ID}|#${GITHUB_RUN_NUMBER}>) of <https://github.com/${GITHUB_REPOSITORY}|${GITHUB_REPOSITORY}> (${LINK}) ${STATUS_MESSAGE}.\n
- ${GITHUB_HEAD_COMMIT_MESSAGE} (<https://github.com/${GITHUB_REPOSITORY}/commit/${GITHUB_SHA}|${COMMIT_HASH}>) by <https://github.com/${GITHUB_REPOSITORY}/pulse|${GITHUB_ACTOR}>
EOF
)

readonly DATA=$(cat <<EOF
{
  "attachments": [
    {
      "blocks": [
        {
          "type": "section",
          "text": {
            "type": "mrkdwn",
            "text": "${TEXT}"
          }
        }
      ],
      "color": "${COLOR}"
    },
  ],
  "blocks": [],
  "channel": "${SLACK_CHANNEL}",
  "text": "",
  "username": "${SLACK_USERNAME}"
}
EOF
)

curl -X POST \
     -H "Content-type: application/json; charset=utf-8" \
     -H "Authorization: Bearer ${SLACK_BOT_USER_OAUTH_TOKEN}" \
     -d "${DATA}" \
     https://slack.com/api/chat.postMessage
