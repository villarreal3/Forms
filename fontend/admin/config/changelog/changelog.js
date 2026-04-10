/**
 * Changelog (Admin UI)
 *
 * Objetivo:
 * - Mantener una lista corta y actual (sin historiales viejos dentro del repo).
 * - Formato simple para que el equipo agregue una entrada nueva en segundos.
 *
 * Cómo agregar una entrada:
 * - Duplica el objeto de ejemplo en `CHANGELOG_ENTRIES`.
 * - Define `version` (ej: "0.3.0"), `date` (YYYY-MM-DD) y `title`.
 * - Agrega cambios por categoría en `changes`.
 */

const CHANGELOG_ENTRIES = [
  {
    version: "Actual",
    date: "2026-03-13",
    title: "Cambios recientes",
    changes: {
      added: [],
      changed: [],
      fixed: [],
      removed: [],
      deprecated: [],
      security: [],
    },
  },
];

const CHANGE_TYPES_ORDER = ["added", "changed", "fixed", "removed", "deprecated", "security"];

const CHANGE_TYPE_LABELS = {
  added: "Agregado",
  changed: "Modificado",
  fixed: "Corregido",
  removed: "Eliminado",
  deprecated: "Obsoleto",
  security: "Seguridad",
};

function formatDate(dateString) {
  const d = new Date(dateString);
  if (Number.isNaN(d.getTime())) return dateString;
  return d.toLocaleDateString("es-ES", { year: "numeric", month: "long", day: "numeric" });
}

function normalizeEntries(entries) {
  return (entries || [])
    .slice()
    .sort((a, b) => new Date(b.date || 0) - new Date(a.date || 0))
    .map((e) => ({
      version: e.version || "N/A",
      date: e.date || "",
      title: e.title || "",
      changes: e.changes || {},
    }));
}

function getTotalChanges(changes, filter) {
  const keys = filter && filter !== "all" ? [filter] : CHANGE_TYPES_ORDER;
  return keys.reduce((sum, k) => sum + (Array.isArray(changes?.[k]) ? changes[k].length : 0), 0);
}

function renderEntry(entry, index, filter) {
  const versionId = `version-${index}`;
  const total = getTotalChanges(entry.changes, filter);
  const subtitle = [formatDate(entry.date), total ? `${total} cambios` : "Sin cambios listados"].filter(Boolean).join(" • ");

  const sections = CHANGE_TYPES_ORDER.map((type) => {
    if (filter && filter !== "all" && filter !== type) return "";
    const items = entry.changes?.[type];
    if (!Array.isArray(items) || items.length === 0) return "";
    return `
      <div class="changelog-change-type" data-type="${type}">
        <h4>${CHANGE_TYPE_LABELS[type] || type}</h4>
        <ul class="changelog-changes-list">
          ${items.map((t) => `<li><span>•</span><span>${escapeHtml(String(t))}</span></li>`).join("")}
        </ul>
      </div>
    `;
  }).join("");

  return `
    <div class="changelog-version-card">
      <button
        class="changelog-version-button"
        data-toggle="${versionId}"
        aria-expanded="false"
        aria-controls="${versionId}-content"
        type="button"
      >
        <div class="changelog-version-header">
          <span class="changelog-version-badge">${escapeHtml(entry.version)}</span>
          <div class="changelog-version-info">
            <h3>${escapeHtml(entry.title)}</h3>
            <p>${escapeHtml(subtitle)}</p>
          </div>
        </div>
        <div class="changelog-version-meta">
          ${total ? `<span class="changelog-changes-count">${total} cambios</span>` : ""}
          <span class="changelog-version-toggle">▼</span>
        </div>
      </button>
      <div id="${versionId}-content" class="changelog-version-content" style="display: none;">
        <div class="changelog-changes-container">
          ${sections || `<p class="changelog-empty">No hay cambios para mostrar con este filtro.</p>`}
        </div>
      </div>
    </div>
  `;
}

function escapeHtml(s) {
  return s
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#039;");
}

function toggleCard(versionId, root) {
  const content = root.querySelector(`#${CSS.escape(versionId)}-content`);
  const btn = root.querySelector(`[data-toggle="${CSS.escape(versionId)}"]`);
  if (!content || !btn) return;
  const expanded = btn.getAttribute("aria-expanded") === "true";
  btn.setAttribute("aria-expanded", expanded ? "false" : "true");
  content.style.display = expanded ? "none" : "block";
}

function renderChangelog() {
  const root = document;
  const container = root.getElementById("changelogContainer");
  const meta = root.getElementById("changelogMeta");
  const filterSelect = root.getElementById("changelogFilter");
  if (!container) return;

  const filter = filterSelect ? filterSelect.value : "all";
  const entries = normalizeEntries(CHANGELOG_ENTRIES);

  container.setAttribute("aria-busy", "false");
  if (!entries.length) {
    container.innerHTML = "<p>No hay actualizaciones disponibles.</p>";
    if (meta) meta.textContent = "";
    return;
  }

  const totalEntries = entries.length;
  const totalChanges = entries.reduce((sum, e) => sum + getTotalChanges(e.changes, filter), 0);
  if (meta) meta.textContent = `${totalEntries} versión(es) • ${totalChanges} cambio(s)`;

  container.innerHTML = entries.map((e, i) => renderEntry(e, i, filter)).join("");

  // toggle handler (event delegation)
  container.addEventListener("click", (ev) => {
    const btn = ev.target && ev.target.closest ? ev.target.closest("[data-toggle]") : null;
    if (!btn) return;
    const versionId = btn.getAttribute("data-toggle");
    if (!versionId) return;
    toggleCard(versionId, root);
  }, { once: true });
}

function initChangelog() {
  const filterSelect = document.getElementById("changelogFilter");
  if (filterSelect) {
    filterSelect.addEventListener("change", () => renderChangelog());
  }
  renderChangelog();
}

// bootstrap.js se carga antes; esperamos el siguiente tick del event loop.
setTimeout(initChangelog, 0);
