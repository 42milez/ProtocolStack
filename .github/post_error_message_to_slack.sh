#!/bin/bash

set -eu

readonly COMMIT_HASH=$(echo "${GITHUB_SHA}" | cut -c 1-7)
readonly LINK=$(
if [[ "${GITHUB_REPOSITORY_REF}" =~ ^"refs/heads/" ]]; then
  readonly BRANCH=$(echo "${GITHUB_REPOSITORY_REF}" | cut -c 12-)
  echo "<https://github.com/${GITHUB_REPOSITORY}/tree/${BRANCH}|${BRANCH}>"
else
  readonly TAG=$(echo "${GITHUB_REPOSITORY_REF}" | cut -c 11-)
  echo "<https://github.com/${GITHUB_REPOSITORY}/releases/tag/${TAG}|${TAG}>"
fi
)
readonly TEXT=$(cat <<EOF
ERROR: ${GITHUB_WORKFLOW} <https://github.com/${GITHUB_REPOSITORY}/actions/runs/${GITHUB_RUN_ID}|#${GITHUB_RUN_NUMBER}> of <https://github.com/${GITHUB_REPOSITORY}|${GITHUB_REPOSITORY}> (${LINK}) failed.\n
- ${GITHUB_HEAD_COMMIT_MESSAGE} (<https://github.com/${GITHUB_REPOSITORY}/commit/${GITHUB_SHA}|${COMMIT_HASH}>) by ${GITHUB_ACTOR}
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
      "color": "danger"
    },
  ],
  "channel": "github",
  "username": "GitHub Support"
}
EOF
)

curl -X POST \
     -H "Content-type: application/json; charset=utf-8" \
     -H "Authorization: Bearer ${SLACK_BOT_USER_OAUTH_TOKEN}" \
     -d "${DATA}" \
     https://slack.com/api/chat.postMessage
