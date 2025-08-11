#!/usr/env bash

PWD=$(pwd)

if [ -z ${INGESTLIST_RELEASE+x} ]; then
	INGESTLIST_RELEASE='7.2.0'
fi

if [ ! -d ${CI_PROJECT_DIR}/buildtools ]; then
  mkdir ${CI_PROJECT_DIR}/buildtools
  pushd ${CI_PROJECT_DIR}/buildtools
  git -c http.sslVerify=false clone https://${GITLAB_USER_LOGIN}:${GITLAB_TOKEN}@gitlab.la-bw.de/dimag/querschnitt/dimagutils.git dimagutils
  mv dimagutils/python-release-download download-release
  rm -rf dimagutils
  cd download-release
  python3 -m venv .venv
  if [ ${OSTYPE} == "msys" ]; then
    source .venv/Scripts/activate
  else
    source .venv/bin/activate
  fi
  # upgrading pip first
  pip install --upgrade pip
  pip install -e .
  popd
fi

if [ ! -d ${CI_PROJECT_DIR}/sources ]; then
  mkdir -p ${CI_PROJECT_DIR}/sources
else 
  rm -rf ${CI_PROJECT_DIR}/sources/*
fi


if [ ! -d ${CI_PROJECT_DIR}/zip ]; then
  mkdir -p ${CI_PROJECT_DIR}/zip
fi

if [ ${OSTYPE} == "msys" ]; then
  export PATH=${PATH}:${PWD}/buildtools/download-release/.venv/Scripts
else
  export PATH=${PATH}:${PWD}/buildtools/download-release/.venv/bin
fi

if [ ! -f zip/ingestlist-${INGESTLIST_RELEASE}.zip ]; then
  download-release --token=${GITLAB_TOKEN} --out=${CI_PROJECT_DIR}/zip/ingestlist-${INGESTLIST_RELEASE}.zip --project=dimag/ingest/ingestlist --release=${INGESTLIST_RELEASE}
fi

unzip -q -o ${CI_PROJECT_DIR}/zip/ingestlist-${INGESTLIST_RELEASE}.zip -d ${CI_PROJECT_DIR}/sources/ingestlist-${INGESTLIST_RELEASE}/
