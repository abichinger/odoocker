ARG TAG
from odoo:${TAG}

USER root

ARG TOUR
RUN if [[ -n "$TOUR" ]]; then pip install websocket-client; fi
RUN if [[ -n "$TOUR" ]]; then \ 
    apt-get update \
    && curl -LO https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb \
    && apt-get install -y ./google-chrome-stable_current_amd64.deb \
    && rm google-chrome-stable_current_amd64.deb; fi

USER odoo

ENTRYPOINT ["/entrypoint.sh"]
CMD ["odoo"]