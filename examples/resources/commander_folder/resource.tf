# commander_folder — create and manage a Keeper vault folder (non-shared) via Commander.
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
#   color           — Folder color: none, red, green, blue, orange, yellow, gray.
#   records         — Set of record UIDs or titles to link into this folder.

# Usage 1 — Basic folder at vault root
resource "commander_folder" "basic" {
  name = "My Documents"
}

# Usage 2 — Folder with a color
resource "commander_folder" "colored" {
  name  = "Important Notes"
  color = "blue"
}

# Usage 3 — Folder inside a parent folder with color and records
resource "commander_folder" "full" {
  name            = "Project Files"
  folder_location = "Work/Engineering"
  color           = "orange"

  records = [
    "QVBa5CNKgyNG_1t_NXzRKg",
    "VUvsSeHTLH7mLnmbok-szg",
  ]
}

# Usage 4 — Minimal folder (only required attribute)
# resource "commander_folder" "minimal" {
#   name = "Quick Notes"
# }

# Usage 5 — Update example: rename, change color, modify records
# resource "commander_folder" "full" {
#   name            = "Project Files Renamed"
#   folder_location = "Work/Engineering"
#   color           = "green"
#
#   records = [
#     "QVBa5CNKgyNG_1t_NXzRKg",
#   ]
# }

# Usage 6 — Move folder to a different parent
# resource "commander_folder" "full" {
#   name            = "Project Files"
#   folder_location = "Archive/2026"
#   color           = "gray"
# }

# Usage 7 — Import (see import.sh)
# resource "commander_folder" "imported" {
# }
# terraform import commander_folder.imported <folder_uid>

output "folder_id" {
  description = "Folder UID for lookups, import, and Commander APIs."
  value       = commander_folder.basic.id
}

output "folder_name" {
  value = commander_folder.basic.name
}

output "full_folder_id" {
  value = commander_folder.full.id
}
