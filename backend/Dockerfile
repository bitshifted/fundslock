# Copyright 2024 Bitshift D.O.O
# SPDX-License-Identifier: MPL-2.0

FROM ubuntu:24.04

ARG BUILD_DATE
ARG VCS_REF
ARG BUILD_VERSION

# Labels.
LABEL org.label-schema.schema-version="1.0"
LABEL org.label-schema.build-date=$BUILD_DATE
LABEL org.label-schema.name="bitshifted/fundslock-be"
LABEL org.label-schema.description="REST API backend for FundsLock smart contract"
LABEL org.label-schema.url="https://github.com/bitshifted/fundslock"
LABEL org.label-schema.vcs-url="https://github.com/bitshifted/fundslock"
LABEL org.label-schema.vcs-ref=$VCS_REF
LABEL org.label-schema.vendor="WSO2"
LABEL org.label-schema.version=$BUILD_VERSION


ARG binary_location=target/linux-amd64/fundslock-be


COPY ${binary_location} /usr/bin
RUN chmod 755 /usr/bin/fundslock-be

ENTRYPOINT [ "/usr/bin/fundslock-be" ]
CMD ["start"]
