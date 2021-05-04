#!/bin/bash

set -eu

readonly COLOR=$(
case ${NEEDS_PREVIOUS_JOB_RESULT} in
  "success") echo "#74c7b8";;
  "failure") echo "#ef4f4f";;
  "cancelled") echo "#f4d160";;
  "skipped") echo "#dddddd";;
  *) {
    echo "Unknown job status."
    exit 1
  }
esac
)

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

readonly STATUS_MESSAGE=$(
case ${NEEDS_PREVIOUS_JOB_RESULT} in
  "success") echo "passed";;
  "failure") echo "failed";;
  "cancelled") echo "was cancelled";;
  "skipped") echo "was skipped";;
  *) {
    echo "Unknown job status."
    exit 1
  }
esac
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
