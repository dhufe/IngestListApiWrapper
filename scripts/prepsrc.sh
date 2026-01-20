#!/usr/bin/env bash

set -euo pipefail  # Exit on error, undefined variables, and pipe failures

# ============================================================================
# Configuration
# ============================================================================

readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly PROJECT_DIR="${CI_PROJECT_DIR:-${SCRIPT_DIR}}"
readonly INGESTLIST_RELEASE="${INGESTLIST_RELEASE:-7.2.0}"

readonly BUILDTOOLS_DIR="${PROJECT_DIR}/buildtools"
readonly DOWNLOAD_RELEASE_DIR="${BUILDTOOLS_DIR}/download-release"
readonly SOURCES_DIR="${PROJECT_DIR}/sources"
readonly ZIP_DIR="${PROJECT_DIR}/zip"

readonly GITLAB_URL="https://gitlab.la-bw.de"
readonly REPO_PATH="dimag/querschnitt/dimagutils.git"

# ============================================================================
# Helper Functions
# ============================================================================

log_info() {
    echo "[INFO] $*" >&2
}

log_error() {
    echo "[ERROR] $*" >&2
}

detect_venv_path() {
    if [[ "${OSTYPE}" == "msys" ]]; then
        echo ".venv/Scripts"
    else
        echo ".venv/bin"
    fi
}

ensure_directory() {
    local dir="$1"
    local clean="${2:-false}"

    if [[ "${clean}" == "true" ]] && [[ -d "${dir}" ]]; then
        log_info "Cleaning directory: ${dir}"
        rm -rf "${dir:?}"/*
    fi

    if [[ ! -d "${dir}" ]]; then
        log_info "Creating directory: ${dir}"
        mkdir -p "${dir}"
    fi
}

# ============================================================================
# Setup Functions
# ============================================================================

setup_buildtools() {
    if [[ -d "${BUILDTOOLS_DIR}" ]]; then
        log_info "Buildtools already exist, skipping setup"
        return 0
    fi

    log_info "Setting up buildtools..."

    ensure_directory "${BUILDTOOLS_DIR}"

    local temp_dir="${BUILDTOOLS_DIR}/dimagutils"

    git -c http.sslVerify=false clone \
        "https://${GITLAB_USER_LOGIN}:${GITLAB_TOKEN}@${GITLAB_URL#https://}/${REPO_PATH}" \
        "${temp_dir}"

    mv "${temp_dir}/python-release-download" "${DOWNLOAD_RELEASE_DIR}"
    rm -rf "${temp_dir}"

    setup_python_venv
}

setup_python_venv() {
    log_info "Setting up Python virtual environment..."

    cd "${DOWNLOAD_RELEASE_DIR}"

    python3 -m venv .venv

    local venv_path
    venv_path="$(detect_venv_path)"

    # shellcheck source=/dev/null
    source "${venv_path}/activate"

    log_info "Upgrading pip..."
    pip install --upgrade pip --quiet

    log_info "Installing package..."
    pip install -e . --quiet

    cd - > /dev/null
}

download_ingestlist() {
    local zip_file="${ZIP_DIR}/ingestlist-${INGESTLIST_RELEASE}.zip"

    if [[ -f "${zip_file}" ]]; then
        log_info "Ingestlist ${INGESTLIST_RELEASE} already downloaded"
        return 0
    fi

    log_info "Downloading ingestlist ${INGESTLIST_RELEASE}..."

    download-release \
        --token="${GITLAB_TOKEN}" \
        --out="${zip_file}" \
        --project=dimag/ingest/ingestlist \
        --release="${INGESTLIST_RELEASE}"
}

extract_ingestlist() {
    local zip_file="${ZIP_DIR}/ingestlist-${INGESTLIST_RELEASE}.zip"
    local extract_dir="${SOURCES_DIR}/ingestlist-${INGESTLIST_RELEASE}"

    log_info "Extracting ingestlist to ${extract_dir}..."

    ensure_directory "${extract_dir}"

    unzip -q -o "${zip_file}" -d "${extract_dir}/"
}

setup_path() {
    local venv_path
    venv_path="$(detect_venv_path)"

    export PATH="${PATH}:${DOWNLOAD_RELEASE_DIR}/${venv_path}"

    log_info "PATH updated with virtual environment"
}

# ============================================================================
# Main
# ============================================================================

main() {
    log_info "Starting setup process..."

    # Validate required environment variables
    if [[ -z "${GITLAB_USER_LOGIN:-}" ]] || [[ -z "${GITLAB_TOKEN:-}" ]]; then
        log_error "GITLAB_USER_LOGIN and GITLAB_TOKEN must be set"
        exit 1
    fi

    # Setup directories
    ensure_directory "${SOURCES_DIR}" "true"
    ensure_directory "${ZIP_DIR}"

    # Setup buildtools and Python environment
    setup_buildtools

    # Update PATH
    setup_path

    # Download and extract ingestlist
    download_ingestlist
    extract_ingestlist

    log_info "Setup completed successfully!"
}

main "$@"
