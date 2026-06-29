import { useMemo, useState } from "react";
import {
  ArrowRight,
  CheckCircle2,
  LogIn,
  LogOut,
  Pencil,
  Plus,
  Save,
  Search,
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
 * The app starts on a presentation page and does not contact the backend until
 * the user submits login or registration. This keeps backend logs clean and
 * makes the first visit feel intentional.
 */
export function App() {
  const api = useMemo(() => createApiClient(), []);
  const [view, setView] = useState("landing");
  const [authMode, setAuthMode] = useState("login");
  const [editingId, setEditingId] = useState(null);
  const [isSaving, setIsSaving] = useState(false);
  const [message, setMessage] = useState("");
  const [notes, setNotes] = useState([]);

  async function runSavingAction(action) {
    setIsSaving(true);

    try {
      await action();
    } catch (error) {
      setMessage(getUserErrorMessage(error));
    } finally {
      setIsSaving(false);
    }
  }

  async function refreshNotes() {
    setNotes(await api.getNotes());
  }

  function openAuth(nextMode = "login") {
    setAuthMode(nextMode);
    setMessage("");
    setView("auth");
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

    await runSavingAction(async () => {
      if (authMode === "register") {
        await api.register(login, password);
      } else {
        await api.login(login, password);
      }

      await refreshNotes();
      setView("notes");
      setMessage("");
    });
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
      await refreshNotes();
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
      setEditingId(null);
      setMessage("");
    });
  }

  async function handleDeleteNote(id) {
    await runSavingAction(async () => {
      await api.deleteNote(id);
      await refreshNotes();
      setMessage("");
    });
  }

  async function handleLogout() {
    await runSavingAction(async () => {
      await api.logout();
      setView("landing");
      setNotes([]);
      setEditingId(null);
      setMessage("");
    });
  }

  if (view === "landing") {
    return <LandingView onLogin={() => openAuth("login")} onRegister={() => openAuth("register")} />;
  }

  if (view === "auth") {
    return (
      <AuthView
        authMode={authMode}
        isSaving={isSaving}
        message={message}
        onBack={() => {
          setMessage("");
          setView("landing");
        }}
        onModeChange={(nextMode) => {
          setAuthMode(nextMode);
          setMessage("");
        }}
        onSubmit={handleAuthSubmit}
      />
    );
  }

  return (
    <NotesView
      editingId={editingId}
      isSaving={isSaving}
      message={message}
      notes={notes}
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

function LandingView({ onLogin, onRegister }) {
  return (
    <main className="landing-layout">
      <section className="landing-hero" aria-label="Описание сервиса">
        <div className="hero-copy">
          <p className="kicker">Личные заметки и задачи</p>
          <h1>Соберите мысли в спокойный рабочий список.</h1>
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
        <Feature icon={<Search size={20} aria-hidden="true" />} title="Видите главное" text="Заметки разложены карточками, чтобы список легко читался." />
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
  message,
  notes,
  onAddNote,
  onCancelEdit,
  onDeleteNote,
  onEditNote,
  onLogout,
  onSaveNote,
}) {
  const count = notes.length;
  const latestNote = notes[0]?.text || "Новая заметка появится здесь после сохранения.";

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

      <section className="workspace-grid">
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

        <aside className="ideas-panel" aria-label="Идеи для заметок">
          <h2>Что можно записать</h2>
          <ul>
            <li>одну главную задачу на сегодня;</li>
            <li>идею, к которой нужно вернуться;</li>
            <li>маленькое дело, которое легко забыть.</li>
          </ul>
        </aside>
      </section>

      <Message text={message} />

      <section className="notes-section" aria-label="Сохраненные заметки">
        <div className="section-heading">
          <h2>{count ? formatNotesCount(count) : "Пока нет заметок"}</h2>
        </div>
        <div className="notes-grid">
          {count ? (
            notes.map((note) => (
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

  return messages.get(error.message) || error.message || "Что-то пошло не так.";
}

function Message({ text }) {
  return text ? (
    <p className="message" role="alert">
      {text}
    </p>
  ) : null;
}
