function toggleAdmin(id) {
  fetch(`/groups/user/${id}/toggle`, { method: "POST" })
    .then(() => location.reload())
    .catch(err => alert("Failed to toggle admin"));
}

function removeMember(id) {
  if (!confirm("Remove member from group?")) return;
  fetch(`/groups/user/${id}/remove`, { method: "DELETE" })
    .then(() => location.reload())
    .catch(err => alert("Failed to remove member"));
}
