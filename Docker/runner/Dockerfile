FROM golang:alpine AS build

RUN apk update && apk add make && apk add git
WORKDIR /
# current directory must be project root evenstore-service-go
COPY . ./

RUN mkdir -p /dist/linux_amd64/release
RUN make build-runner

# Image to use
FROM ubuntu:18.04 AS runtime

# Powershell
ARG PS_VERSION=7.0.3
ARG PS_PACKAGE=powershell-lts_${PS_VERSION}-1.ubuntu.18.04_amd64.deb
ARG PS_PACKAGE_URL=https://github.com/PowerShell/PowerShell/releases/download/v${PS_VERSION}/${PS_PACKAGE}

# Define ENVs for Localization/Globalization
ENV DOTNET_SYSTEM_GLOBALIZATION_INVARIANT=false \
    LC_ALL=en_US.UTF-8 \
    LANG=en_US.UTF-8 \
    # set a fixed location for the Module analysis cache
    PSModuleAnalysisCachePath=/var/cache/microsoft/powershell/PSModuleAnalysisCache/ModuleAnalysisCache \
    POWERSHELL_DISTRIBUTION_CHANNEL=PSDocker-Ubuntu-18.04

# Install dependencies and clean up
RUN apt-get update \
    && apt-get install --no-install-recommends -y \
    # curl is required to grab the Linux package
        curl \
    # requied to setup the locale
        locales \
    # required for SSL
        ca-certificates \
        gss-ntlmssp \
        wget \
        unzip \
    # Download the Linux package and save it
    && echo ${PS_PACKAGE_URL} \
    && curl -sSL ${PS_PACKAGE_URL} -o /tmp/powershell.deb \
    && apt-get install --no-install-recommends -y /tmp/powershell.deb \
    && apt-get dist-upgrade -y \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* \
    && locale-gen $LANG && update-locale \
    # remove powershell package
    && rm /tmp/powershell.deb \
    # intialize powershell module cache
    # and disable telemetry
    && export POWERSHELL_TELEMETRY_OPTOUT=1 \
    && pwsh \
        -NoLogo \
        -NoProfile \
        -Command " \
          \$ErrorActionPreference = 'Stop' ; \
          \$ProgressPreference = 'SilentlyContinue' ; \
          while(!(Test-Path -Path \$env:PSModuleAnalysisCachePath)) {  \
            Write-Host "'Waiting for $env:PSModuleAnalysisCachePath'" ; \
            Start-Sleep -Seconds 6 ; \
          }"

# Install Az module
RUN pwsh -Command "Install-Module -Name Az -AllowClobber -Scope AllUsers -Repository PSGallery -Force"

# Create user
RUN adduser --disabled-password --gecos "" --uid 10001 appuser

# create artifacts directory
RUN mkdir artifacts
RUN chown -R appuser /artifacts


WORKDIR /
COPY --from=build /dist/linux_amd64/release /
COPY --from=build /cmd/runner/azpwsh.ps1 /

RUN chown appuser ./azpwsh.ps1

USER appuser:appuser

CMD [ "pwsh" ]