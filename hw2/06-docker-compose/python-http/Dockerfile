FROM python:3.13-slim-trixie

RUN apt-get update && \
  apt-get install -y --no-install-recommends curl && \
  rm -rf /var/lib/apt/lists/*

WORKDIR /app

RUN groupadd -g 1000 appuser && \
  useradd -u 1000 -g appuser -m appuser && \
  chown -R appuser:appuser /app

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

USER appuser
COPY --chown=appuser:appuser . .

HEALTHCHECK --interval=30s --timeout=5s CMD curl -f http://localhost:5000/health
CMD ["python", "app.py"]
