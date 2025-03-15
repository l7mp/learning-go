cd "${CONTEXT_PATH}" || exit 1
if [[ -z "${FILE_PATH}" ]]; then
  docker build -t "${TAGS}" .
else
  docker build -f "${FILE_PATH}" -t "${TAGS}" .
fi
IMAGE_NAME=$(echo "${TAGS}" | cut -d : -f1)
mkdir images
docker image save -o "images/${IMAGE_NAME}" "${TAGS}"
