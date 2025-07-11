function updateColorPreview() {
  const input = document.getElementById("colorInput");
  const preview = document.getElementById("colorPreview");

  // Remove any previous color classes
  preview.className = "w-8 h-8 rounded border";

  // Add new class dynamically
  const colorClass = input.value.trim();
  if (colorClass) preview.classList.add(colorClass);
}