package payload

func GetHarborProjects() string {
	return `[
	{
		"chart_count": 0,
		"creation_time": "2021-10-07T11:06:22.000Z",
		"current_user_role_id": 1,
		"current_user_role_ids": [
		1
		],
		"cve_allowlist": {
		"creation_time": "0001-01-01T00:00:00.000Z",
		"id": 2,
		"items": [],
		"project_id": 2,
		"update_time": "0001-01-01T00:00:00.000Z"
		},
		"metadata": {
		"public": "false"
		},
		"name": "dummy",
		"owner_id": 3,
		"owner_name": "user",
		"project_id": 2,
		"repo_count": 1,
		"update_time": "2021-10-07T11:06:22.000Z"
	},
	{
		"chart_count": 0,
		"creation_time": "2022-02-02T08:23:46.000Z",
		"current_user_role_id": 1,
		"current_user_role_ids": [
		1
		],
		"cve_allowlist": {
		"creation_time": "0001-01-01T00:00:00.000Z",
		"id": 7,
		"items": [],
		"project_id": 7,
		"update_time": "0001-01-01T00:00:00.000Z"
		},
		"metadata": {
		"public": "false"
		},
		"name": "web",
		"owner_id": 3,
		"owner_name": "user",
		"project_id": 7,
		"repo_count": 166,
		"update_time": "2022-02-02T08:23:46.000Z"
	}
]`
}

func GetHarborRepos() string {
	return `[
  {
    "artifact_count": 1,
    "creation_time": "2024-01-08T20:02:39.634Z",
    "id": 1983,
    "name": "web/app-feature-1315",
    "project_id": 7,
    "pull_count": 5,
    "update_time": "2024-01-10T11:16:23.979Z"
  },
  {
    "artifact_count": 1,
    "creation_time": "2024-01-08T20:02:37.514Z",
    "id": 1982,
    "name": "web/app",
    "project_id": 7,
    "pull_count": 4,
    "update_time": "2024-01-10T05:02:11.314Z"
  },
  {
    "artifact_count": 1,
    "creation_time": "2024-01-08T20:02:21.586Z",
    "id": 1981,
    "name": "web/app-feature-1315-builder",
    "project_id": 7,
    "pull_count": 2,
    "update_time": "2024-01-08T20:02:49.997Z"
  }
]`
}

func GetHarborArtifacts() string {
	return `[
  {
    "accessories": null,
    "addition_links": {
      "build_history": {
        "absolute": false,
        "href": "/api/v2.0/projects/web/repositories/app/artifacts/sha256:xxxxxxxxxxxxxx/additions/build_history"
      }
    },
    "digest": "sha256:xxxxxxxxxxxxxx",
    "id": 29422,
    "labels": null,
    "project_id": 7,
    "pull_time": "2024-01-10T12:34:50.675Z",
    "push_time": "2024-01-10T11:36:06.917Z",
    "references": null,
    "repository_id": 407,
    "size": 42384496,
    "tags": [
      {
        "artifact_id": 29422,
        "id": 18447,
        "immutable": false,
        "name": "0_0_xxxxx",
        "pull_time": "0001-01-01T00:00:00.000Z",
        "push_time": "2024-01-10T11:36:07.881Z",
        "repository_id": 407,
        "signed": false
      }
    ],
    "type": "IMAGE"
  },
  {
    "accessories": null,
    "addition_links": {
      "build_history": {
        "absolute": false,
        "href": "/api/v2.0/projects/web/repositories/app/artifacts/sha256:yyyyyyyyyyyyyyyy/additions/build_history"
      }
    },
    "digest": "sha256:yyyyyyyyyyyyyyyy",
    "id": 29321,
    "labels": null,
    "project_id": 7,
    "pull_time": "2024-01-09T10:18:51.849Z",
    "push_time": "2024-01-09T10:18:04.059Z",
    "references": null,
    "repository_id": 407,
    "size": 42384426,
    "tags": [
      {
        "artifact_id": 29321,
        "id": 18360,
        "immutable": false,
        "name": "0_0_yyyyyy",
        "pull_time": "0001-01-01T00:00:00.000Z",
        "push_time": "2024-01-09T10:18:04.343Z",
        "repository_id": 407,
        "signed": false
      }
    ],
    "type": "IMAGE"
  },
  {
    "accessories": null,
    "addition_links": {
      "build_history": {
        "absolute": false,
        "href": "/api/v2.0/projects/web/repositories/app/artifacts/sha256:zzzzzzzzzzzzzzz/additions/build_history"
      }
    },
    "digest": "sha256:zzzzzzzzzzzzzzz",
    "id": 29220,
    "labels": null,
    "project_id": 7,
    "pull_time": "2024-01-09T05:00:30.624Z",
    "push_time": "2024-01-08T12:51:07.966Z",
    "references": null,
    "repository_id": 407,
    "size": 42384427,
    "tags": [
      {
        "artifact_id": 29220,
        "id": 18301,
        "immutable": false,
        "name": "0_0_zzzzzz",
        "pull_time": "0001-01-01T00:00:00.000Z",
        "push_time": "2024-01-08T12:51:08.326Z",
        "repository_id": 407,
        "signed": false
      }
    ],
    "type": "IMAGE"
  }
]`
}
