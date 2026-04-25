"""Import WSY API definitions into Apifox."""
import json
import os
import sys
import requests

# Load .env
script_dir = os.path.dirname(os.path.abspath(__file__))
env_path = os.path.join(script_dir, "..", "..", "apifox-mcp", ".env")

APIFOX_TOKEN = None
APIFOX_PROJECT_ID = None

with open(env_path) as f:
    for line in f:
        line = line.strip()
        if line.startswith("APIFOX_TOKEN="):
            APIFOX_TOKEN = line.split("=", 1)[1]
        elif line.startswith("APIFOX_PROJECT_ID="):
            APIFOX_PROJECT_ID = line.split("=", 1)[1]

if not APIFOX_TOKEN or not APIFOX_PROJECT_ID:
    print("ERROR: Missing APIFOX_TOKEN or APIFOX_PROJECT_ID")
    sys.exit(1)

# Read OpenAPI spec
spec_path = os.path.join(script_dir, "wsy-api.json")
with open(spec_path, "r") as f:
    openapi_spec = json.dumps(json.load(f), ensure_ascii=False)

# Import to Apifox
url = f"https://api.apifox.com/v1/projects/{APIFOX_PROJECT_ID}/import-openapi"
headers = {
    "Authorization": f"Bearer {APIFOX_TOKEN}",
    "Content-Type": "application/json",
    "X-Apifox-Api-Version": "2024-03-28",
}

payload = {
    "input": openapi_spec,
    "options": {
        "targetEndpointFolderId": 0,  # root level, will create folder automatically
        "targetSchemaFolderId": 0,
        "endpointOverwriteBehavior": "CREATE_NEW",
        "schemaOverwriteBehavior": "OVERWRITE_EXISTING",
    },
}

print(f"Importing WSY API to Apifox project {APIFOX_PROJECT_ID}...")
resp = requests.post(url, json=payload, params={"locale": "zh-CN"}, headers=headers)
result = resp.json()

print(f"Status: {resp.status_code}")
print(json.dumps(result, indent=2, ensure_ascii=False))
