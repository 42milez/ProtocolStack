#!/bin/bash

set -eu

readonly TEXT=$(cat <<EOF
${GITHUB_WORKFLOW} <https://github.com/${GITHUB_REPOSITORY}/actions/runs/${GITHUB_RUN_ID}|#${GITHUB_RUN_NUMBER}> of <https://github.com/${GITHUB_REPOSITORY}|${GITHUB_REPOSITORY}> (${GITHUB_REPOSITORY_REF}) failed.\n
- ${GITHUB_HEAD_COMMIT_MESSAGE} (<https://github.com/${GITHUB_REPOSITORY}/commit/${GITHUB_SHA}|$(echo "${GITHUB_SHA}" | cut -c 7 )>) by ${GITHUB_ACTOR}
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
    },
    "color": "danger"
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
