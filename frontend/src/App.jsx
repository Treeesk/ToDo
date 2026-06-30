import { useEffect, useMemo, useState } from "react";
import {
  ArrowRight,
  CheckCircle2,
  LogIn,
  LogOut,
  Pencil,
  Plus,
  Save,
  ShieldCheck,
  Sparkles,
  Trash2,
  UserPlus,
  X,
} from "lucide-react";
import { ApiError, createApiClient } from "./api.js";

/**
 * Root application component.
 *
 * The current URL is the source of truth:
 * "/" is public, "/login" and "/register" submit auth requests, and
 * "/notes" is protected by loading notes from the backend with HttpOnly
 * cookies.
 */
export function App() {
  const api = useMemo(() => createApiClient(), []);
  const [path, setPath] = useState(getCurrentPath);
  const [editingId, setEditingId] = useState(null);
  const [isLoadingNotes, setIsLoadingNotes] = useState(false);
  const [isSaving, setIsSaving] = useState(false);
  const [message, setMessage] = useState("");
  const [notes, setNotes] = useState([]);
  const [latestNoteText, setLatestNoteText] = useState("");
  const [recentNoteId, setRecentNoteId] = useState(null);

  useEffect(() => {
    function handlePopState() {
      setPath(getCurrentPath());
      setMessage("");
    }

    window.addEventListener("popstate", handlePopState);

    return () => {
      window.removeEventListener("popstate", handlePopState);
    };
  }, []);

  useEffect(() => {
    if (path !== "/notes") {
      return;
    }

    let isMounted = true;
    setIsLoadingNotes(true);

    async function loadProtectedNotes() {
      try {
        const loadedNotes = await api.getNotes();

        if (!isMounted) {
          return;
        }

        applyNotesState(loadedNotes);
        setMessage("");
      } catch (error) {
        if (!isMounted) {
          return;
        }

        handleAuthFailure(error);
      } finally {
        if (isMounted) {
          setIsLoadingNotes(false);
        }
      }
    }

    loadProtectedNotes();

    return () => {
      isMounted = false;
    };
  }, [api, path]);

  async function runSavingAction(action) {
    setIsSaving(true);

    try {
      await action();
    } catch (error) {
      handleAuthFailure(error);
    } finally {
      setIsSaving(false);
    }
  }

  async function refreshNotes() {
    const loadedNotes = await api.getNotes();
    applyNotesState(loadedNotes);
    return loadedNotes;
  }

  function navigate(nextPath, { replace = false } = {}) {
    if (getCurrentPath() !== nextPath) {
      const method = replace ? "replaceState" : "pushState";
      window.history[method](null, "", nextPath);
    }

    setPath(nextPath);
    setMessage("");

    if (nextPath !== "/notes") {
      clearNotesState();
    }
  }

  function applyNotesState(loadedNotes, latestOverride = null, recentOverride = null) {
    const newestNote = getNewestNote(loadedNotes);
    setNotes(loadedNotes);
    setLatestNoteText(latestOverride ?? newestNote?.text ?? "");
    setRecentNoteId(recentOverride ?? newestNote?.id ?? null);
  }

  function clearNotesState() {
    setNotes([]);
    setLatestNoteText("");
    setRecentNoteId(null);
    setEditingId(null);
  }

  function handleAuthFailure(error) {
    if (isUnauthorizedError(error)) {
      clearNotesState();
      navigate("/", { replace: true });
      return;
    }

    setMessage(getUserErrorMessage(error));
  }

  async function handleAuthSubmit(event) {
    event.preventDefault();
    const formData = new FormData(event.currentTarget);
    const login = String(formData.get("login") || "").trim();
    const password = String(formData.get("password") || "");

    if (!login || !password) {
      setMessage("Введите логин и пароль.");
      return;
    }

    setIsSaving(true);

    try {
      if (path === "/register") {
        await api.register(login, password);
      } else {
        await api.login(login, password);
      }

      navigate("/notes", { replace: true });
    } catch (error) {
      setMessage(getUserErrorMessage(error));
    } finally {
      setIsSaving(false);
    }
  }

  async function handleAddNote(event) {
    event.preventDefault();
    const formElement = event.currentTarget;
    const formData = new FormData(formElement);
    const text = String(formData.get("text") || "").trim();

    if (!text) {
      setMessage("Сначала напишите заметку.");
      return;
    }

    await runSavingAction(async () => {
      await api.addNote(text);
      formElement.reset();
      const loadedNotes = await refreshNotes();
      const newestNote = getNewestNote(loadedNotes);
      setLatestNoteText(newestNote?.text || text);
      setRecentNoteId(newestNote?.id || null);
      setMessage("");
    });
  }

  async function handleEditNote(event) {
    event.preventDefault();
    const formData = new FormData(event.currentTarget);
    const id = Number(formData.get("id"));
    const text = String(formData.get("text") || "").trim();

    if (!text) {
      setMessage("Текст заметки не может быть пустым.");
      return;
    }

    await runSavingAction(async () => {
      await api.editNote(id, text);
      await refreshNotes();
      setLatestNoteText(text);
      setRecentNoteId(id);
      setEditingId(null);
      setMessage("");
    });
  }

  async function handleDeleteNote(id) {
    await runSavingAction(async () => {
      await api.deleteNote(id);
      const loadedNotes = await refreshNotes();
      const newestNote = getNewestNote(loadedNotes);
      setLatestNoteText(newestNote?.text || "");
      setRecentNoteId(newestNote?.id || null);
      setMessage("");
    });
  }

  async function handleLogout() {
    await runSavingAction(async () => {
      await api.logout();
      clearNotesState();
      navigate("/", { replace: true });
    });
  }

  if (path === "/login" || path === "/register") {
    return (
      <AuthView
        authMode={path === "/register" ? "register" : "login"}
        isSaving={isSaving}
        message={message}
        onBack={() => navigate("/")}
        onModeChange={(nextMode) => navigate(nextMode === "register" ? "/register" : "/login", { replace: true })}
        onSubmit={handleAuthSubmit}
      />
    );
  }

  if (path === "/notes") {
    if (isLoadingNotes) {
      return <LoadingView />;
    }

    return (
      <NotesView
        editingId={editingId}
        isSaving={isSaving}
        latestNoteText={latestNoteText}
        message={message}
        notes={notes}
        recentNoteId={recentNoteId}
        onAddNote={handleAddNote}
        onCancelEdit={() => setEditingId(null)}
        onDeleteNote={handleDeleteNote}
        onEditNote={(id) => {
          setEditingId(id);
          setMessage("");
        }}
        onLogout={handleLogout}
        onSaveNote={handleEditNote}
      />
    );
  }

  return (
    <LandingView
      onLogin={() => navigate("/login")}
      onRegister={() => navigate("/register")}
    />
  );
}

function LoadingView() {
  return (
    <main className="workspace-shell workspace-shell--center">
      <section className="status-panel" aria-live="polite">
        <div className="loader" aria-hidden="true" />
        <p>Загружаем заметки...</p>
      </section>
    </main>
  );
}

function LandingView({ onLogin, onRegister }) {
  return (
    <main className="landing-layout">
      <section className="landing-hero" aria-label="Описание сервиса">
        <div className="hero-copy">
          <p className="kicker">Личные заметки и задачи</p>
          <h1>Заметки без лишнего шума.</h1>
          <p className="intro">
            Быстро сохраняйте задачи, идеи и напоминания. Возвращайтесь к ним тогда, когда
            действительно нужно действовать.
          </p>
          <div className="hero-actions">
            <button className="primary-action" type="button" onClick={onLogin}>
              <LogIn size={18} aria-hidden="true" />
              Войти
            </button>
            <button className="ghost-action" type="button" onClick={onRegister}>
              <UserPlus size={18} aria-hidden="true" />
              Создать аккаунт
            </button>
          </div>
        </div>

        <div className="preview-board" aria-label="Пример рабочего списка">
          <div className="preview-toolbar">
            <span />
            <span />
            <span />
          </div>
          <div className="preview-note preview-note--strong">
            <CheckCircle2 size={20} aria-hidden="true" />
            <p>Подготовить план на неделю</p>
          </div>
          <div className="preview-note">
            <Sparkles size={20} aria-hidden="true" />
            <p>Записать идею для нового проекта</p>
          </div>
          <div className="preview-note">
            <ShieldCheck size={20} aria-hidden="true" />
            <p>Проверить важные дела вечером</p>
          </div>
        </div>
      </section>

      <section className="landing-features" aria-label="Возможности">
        <Feature icon={<Plus size={20} aria-hidden="true" />} title="Добавляйте быстро" text="Короткая форма всегда под рукой после входа." />
        <Feature icon={<Pencil size={20} aria-hidden="true" />} title="Исправляйте без лишних шагов" text="Редактирование происходит прямо в карточке заметки." />
        <Feature icon={<ShieldCheck size={20} aria-hidden="true" />} title="Видите главное" text="Заметки разложены карточками, чтобы список легко читался." />
      </section>
    </main>
  );
}

function Feature({ icon, title, text }) {
  return (
    <article className="feature-item">
      <div className="feature-icon">{icon}</div>
      <h2>{title}</h2>
      <p>{text}</p>
    </article>
  );
}

function AuthView({ authMode, isSaving, message, onBack, onModeChange, onSubmit }) {
  const isRegister = authMode === "register";

  return (
    <main className="auth-layout">
      <section className="brand-panel" aria-label="О сервисе">
        <p className="kicker">{isRegister ? "Новый аккаунт" : "Вход в аккаунт"}</p>
        <h1>{isRegister ? "Создайте свое место для заметок." : "Вернитесь к своим заметкам."}</h1>
        <p className="intro">
          Сохраняйте мысли сразу, пока они не растворились в потоке дел. Остальное интерфейс
          оставит аккуратным и понятным.
        </p>
        <button className="link-action" type="button" onClick={onBack}>
          <ArrowRight size={17} aria-hidden="true" />
          На стартовую
        </button>
      </section>

      <section className="auth-panel" aria-label={isRegister ? "Создать аккаунт" : "Войти"}>
        <div className="segmented-control" role="tablist" aria-label="Режим авторизации">
          <button
            type="button"
            className={!isRegister ? "is-active" : ""}
            aria-selected={!isRegister}
            onClick={() => onModeChange("login")}
          >
            Войти
          </button>
          <button
            type="button"
            className={isRegister ? "is-active" : ""}
            aria-selected={isRegister}
            onClick={() => onModeChange("register")}
          >
            Регистрация
          </button>
        </div>

        <form className="stack" onSubmit={onSubmit}>
          <label>
            <span>Логин</span>
            <input name="login" autoComplete="username" required />
          </label>
          <label>
            <span>Пароль</span>
            <input
              name="password"
              type="password"
              autoComplete={isRegister ? "new-password" : "current-password"}
              required
            />
          </label>
          <Message text={message} />
          <button className="primary-action" type="submit" disabled={isSaving}>
            {isSaving ? (
              "Подождите..."
            ) : isRegister ? (
              <>
                <UserPlus size={18} aria-hidden="true" />
                Создать аккаунт
              </>
            ) : (
              <>
                <LogIn size={18} aria-hidden="true" />
                Войти
              </>
            )}
          </button>
        </form>
      </section>
    </main>
  );
}

function NotesView({
  editingId,
  isSaving,
  latestNoteText,
  message,
  notes,
  recentNoteId,
  onAddNote,
  onCancelEdit,
  onDeleteNote,
  onEditNote,
  onLogout,
  onSaveNote,
}) {
  const orderedNotes = orderNotes(notes, recentNoteId);
  const count = orderedNotes.length;
  const latestNote =
    latestNoteText || orderedNotes[0]?.text || "Новая заметка появится здесь после сохранения.";

  return (
    <main className="workspace-shell">
      <header className="workspace-header">
        <div>
          <p className="kicker">Рабочий список</p>
          <h1>Ваши заметки</h1>
        </div>
        <button className="ghost-action" type="button" onClick={onLogout} disabled={isSaving}>
          <LogOut size={18} aria-hidden="true" />
          Выйти
        </button>
      </header>

      <section className="workspace-summary" aria-label="Сводка по заметкам">
        <article className="summary-tile">
          <span>{count}</span>
          <p>{count ? "заметок сохранено" : "заметок пока нет"}</p>
        </article>
        <article className="summary-tile summary-tile--wide">
          <p>Последняя запись</p>
          <strong>{latestNote}</strong>
        </article>
      </section>

      <section className="composer" aria-label="Добавить заметку">
        <form className="composer-form" onSubmit={onAddNote}>
          <label>
            <span>Новая заметка</span>
            <textarea
              name="text"
              rows="5"
              placeholder="Например: проверить миграции, купить продукты, записать идею..."
              required
            />
          </label>
          <button className="primary-action" type="submit" disabled={isSaving}>
            <Plus size={18} aria-hidden="true" />
            Добавить
          </button>
        </form>
      </section>

      <Message text={message} />

      <section className="notes-section" aria-label="Сохраненные заметки">
        <div className="section-heading">
          <h2>{count ? formatNotesCount(count) : "Пока нет заметок"}</h2>
        </div>
        <div className="notes-grid">
          {count ? (
            orderedNotes.map((note) => (
              <NoteCard
                key={note.id}
                isEditing={editingId === note.id}
                isSaving={isSaving}
                note={note}
                onCancelEdit={onCancelEdit}
                onDeleteNote={onDeleteNote}
                onEditNote={onEditNote}
                onSaveNote={onSaveNote}
              />
            ))
          ) : (
            <EmptyState />
          )}
        </div>
      </section>
    </main>
  );
}

function NoteCard({
  isEditing,
  isSaving,
  note,
  onCancelEdit,
  onDeleteNote,
  onEditNote,
  onSaveNote,
}) {
  if (isEditing) {
    return (
      <article className="note-card">
        <form className="edit-form" onSubmit={onSaveNote}>
          <input type="hidden" name="id" value={note.id} />
          <label>
            <span>Редактировать заметку</span>
            <textarea name="text" rows="5" defaultValue={note.text} required />
          </label>
          <div className="note-actions">
            <button className="primary-action compact" type="submit" disabled={isSaving}>
              <Save size={16} aria-hidden="true" />
              Сохранить
            </button>
            <button className="ghost-action compact" type="button" onClick={onCancelEdit}>
              <X size={16} aria-hidden="true" />
              Отмена
            </button>
          </div>
        </form>
      </article>
    );
  }

  return (
    <article className="note-card">
      <p>{note.text}</p>
      <div className="note-actions">
        <button className="ghost-action compact" type="button" onClick={() => onEditNote(note.id)}>
          <Pencil size={16} aria-hidden="true" />
          Изменить
        </button>
        <button
          className="danger-action compact"
          type="button"
          onClick={() => onDeleteNote(note.id)}
          disabled={isSaving}
        >
          <Trash2 size={16} aria-hidden="true" />
          Удалить
        </button>
      </div>
    </article>
  );
}

function EmptyState() {
  return (
    <div className="empty-state">
      <h3>Список пуст.</h3>
      <p>Добавьте первую заметку выше.</p>
    </div>
  );
}

function formatNotesCount(count) {
  const lastDigit = count % 10;
  const lastTwoDigits = count % 100;

  if (lastDigit === 1 && lastTwoDigits !== 11) {
    return `${count} сохраненная заметка`;
  }

  if (lastDigit >= 2 && lastDigit <= 4 && (lastTwoDigits < 12 || lastTwoDigits > 14)) {
    return `${count} сохраненные заметки`;
  }

  return `${count} сохраненных заметок`;
}

function getUserErrorMessage(error) {
  if (!(error instanceof Error)) {
    return "Что-то пошло не так.";
  }

  const rawMessage = error.message || "";
  const normalizedMessage = rawMessage.toLowerCase();

  if (normalizedMessage.includes("invalid login or password")) {
    return "Неверный логин или пароль.";
  }

  if (normalizedMessage.includes("already exists")) {
    return "Пользователь с таким логином уже есть.";
  }

  const messages = new Map([
    ["missing cookie", "Войдите в аккаунт, чтобы открыть заметки."],
    ["Unauthorized", "Войдите в аккаунт, чтобы продолжить."],
    ["Unauthorized user", "Неверный логин или пароль."],
    ["Error: A user with that login already exists.", "Пользователь с таким логином уже есть."],
    ["text of note is required", "Текст заметки обязателен."],
    ["text is required", "Текст заметки обязателен."],
    ["the text field is missing", "Не найден текст заметки."],
    ["the id field is missing", "Не найден номер заметки."],
    ["Error: input is too long", "Введенный текст слишком длинный."],
    ["Error: body of json empty", "Пустой запрос."],
    ["Error: syntax error", "Некорректный формат запроса."],
    ["Request failed", "Не удалось выполнить запрос."],
  ]);

  return messages.get(rawMessage) || "Не удалось выполнить запрос.";
}

function isUnauthorizedError(error) {
  return error instanceof ApiError && error.status === 401;
}

function getCurrentPath() {
  const allowedPaths = new Set(["/", "/login", "/register", "/notes"]);
  return allowedPaths.has(window.location.pathname) ? window.location.pathname : "/";
}

function getNewestNote(notes) {
  return [...notes].sort((left, right) => right.id - left.id)[0] || null;
}

function orderNotes(notes, recentNoteId) {
  return [...notes].sort((left, right) => {
    if (left.id === recentNoteId) {
      return -1;
    }

    if (right.id === recentNoteId) {
      return 1;
    }

    return right.id - left.id;
  });
}

function Message({ text }) {
  return text ? (
    <p className="message" role="alert">
      {text}
    </p>
  ) : null;
}
