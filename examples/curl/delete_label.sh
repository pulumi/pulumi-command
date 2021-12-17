curl \
  -s \
  -X DELETE \
  -H "authorization: Bearer $GITHUB_TOKEN" \
  -H "Accept: application/vnd.github.v3+json" \
  https://api.github.com/repos/$OWNER/$REPO/labels/$NAME