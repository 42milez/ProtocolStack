#!/bin/bash

set -eu

readonly DATA=$(cat <<EOF
{
  "attachments": [
    {
      "blocks": [
        {
          "type": "section",
          "text": {
            "type": "mrkdwn",
            "text": "${WORKFLOW_NAME} <https://github.com/${REPOSITORY_OWNER}/${REPOSITORY_NAME}/actions/runs/${WORKFLOW_RUN_ID}|#${WORKFLOW_RUN_NUMBER} (${WORKFLOW_RUN_ID})> of ${REPOSITORY_OWNER}/${REPOSITORY_NAME}@${REPOSITORY_REF} by ${ACTOR} failed."
          }
        }
      ],
      "color": "danger"
    }
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
