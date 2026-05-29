# commander_non_shared_folder — create and manage a Keeper vault folder (non-shared) via Commander.
#
# Required:
#   name — Folder name (leaf name, without parent path).
#
# Computed:
#   id — Folder UID (stable; use for import and data source lookups).
#
# Optional:
#   folder_location — Parent folder path where the folder is created (e.g. "Parent/Child").
#                     Leave empty or omit for vault root.
#   records         — Set of record UIDs or titles to link into this folder.

# Usage 1 — Basic folder at vault root
resource "commander_non_shared_folder" "basic" {
  name = "My Documents"
}

# Usage 2 — Folder inside a parent folder with records
resource "commander_non_shared_folder" "full" {
  name            = "Project Files"
  folder_location = "Work/Engineering"

  records = [
    "QVBa5CNKgyNG_1t_NXzRKg",
    "VUvsSeHTLH7mLnmbok-szg",
  ]
}

# Usage 3 — Minimal folder (only required attribute)
# resource "commander_non_shared_folder" "minimal" {
#   name = "Quick Notes"
# }

# Usage 4 — Update example: rename, modify records
# resource "commander_non_shared_folder" "full" {
#   name            = "Project Files Renamed"
#   folder_location = "Work/Engineering"
#
#   records = [
#     "QVBa5CNKgyNG_1t_NXzRKg",
#   ]
# }

# Usage 5 — Move folder to a different parent
# resource "commander_non_shared_folder" "full" {
#   name            = "Project Files"
#   folder_location = "Archive/2026"
# }

# Usage 6 — Import (see import.sh)
# resource "commander_non_shared_folder" "imported" {
# }
# terraform import commander_non_shared_folder.imported <folder_uid>

output "folder_id" {
  description = "Folder UID for lookups, import, and Commander APIs."
  value       = commander_non_shared_folder.basic.id
}

output "folder_name" {
  value = commander_non_shared_folder.basic.name
}

output "full_folder_id" {
  value = commander_non_shared_folder.full.id
}
