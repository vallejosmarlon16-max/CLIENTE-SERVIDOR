import React, { useState, useEffect } from 'react';
import { Plus, ListTodo, Loader2, AlertCircle } from 'lucide-react';
import TaskCard from './components/TaskCard';

const API_URL = 'http://localhost:8081/api/tasks';

export default function App() {
  const [tasks, setTasks] = useState([]);
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // Obtener tareas del backend al montar el componente
  useEffect(() => {
    fetchTasks();
  }, []);

  const fetchTasks = async () => {
    try {
      setLoading(true);
      setError(null);
      const res = await fetch(API_URL);
      if (!res.ok) throw new Error('Error al obtener las tareas');
      const data = await res.json();
      setTasks(data);
    } catch (err) {
      console.error(err);
      setError('No se pudo conectar con el servidor backend.');
    } finally {
      setLoading(false);
    }
  };

  const handleCreateTask = async (e) => {
    e.preventDefault();
    if (!title.trim()) return;

    try {
      const res = await fetch(API_URL, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ title, description }),
      });
      
      if (!res.ok) throw new Error('Error al crear la tarea');
      const newTask = await res.json();
      setTasks((prev) => [newTask, ...prev]);
      setTitle('');
      setDescription('');
    } catch (err) {
      console.error(err);
      alert('Error al guardar la tarea');
    }
  };

  const handleToggleTask = async (task) => {
    try {
      const updatedStatus = !task.completed;
      const res = await fetch(`${API_URL}?id=${task.id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ completed: updatedStatus }),
      });

      if (!res.ok) throw new Error('Error al actualizar la tarea');
      
      setTasks((prev) =>
        prev.map((t) => (t.id === task.id ? { ...t, completed: updatedStatus } : t))
      );
    } catch (err) {
      console.error(err);
      alert('Error al actualizar el estado de la tarea');
    }
  };

  const handleDeleteTask = async (id) => {
    if (!confirm('¿Estás seguro de que deseas eliminar esta tarea?')) return;
    try {
      const res = await fetch(`${API_URL}?id=${id}`, {
        method: 'DELETE',
      });

      if (!res.ok) throw new Error('Error al eliminar la tarea');
      
      setTasks((prev) => prev.filter((t) => t.id !== id));
    } catch (err) {
      console.error(err);
      alert('Error al eliminar la tarea');
    }
  };

  const completedCount = tasks.filter((t) => t.completed).length;

  return (
    <div className="container">
      <header>
        <div className="logo-section">
          <div className="logo-icon">
            <ListTodo size={24} color="#ffffff" />
          </div>
          <div>
            <h1>TaskSphere</h1>
            <p className="subtitle">Gestión inteligente de tus proyectos cotidianos</p>
          </div>
        </div>
        <div className="task-counter">
          <span className="badge badge-pending" style={{ marginRight: '8px' }}>
            Pendientes: {tasks.length - completedCount}
          </span>
          <span className="badge badge-completed">
            Completadas: {completedCount}
          </span>
        </div>
      </header>

      <div className="glass-card form-box">
        <h2 className="form-title">
          <Plus size={20} color="#6366f1" />
          Crear Nueva Tarea
        </h2>
        <form onSubmit={handleCreateTask}>
          <div className="form-group">
            <label htmlFor="title">Título de la Tarea</label>
            <input
              type="text"
              id="title"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              className="form-control"
              placeholder="Ej. Diseñar arquitectura de base de datos"
              required
            />
          </div>
          <div className="form-group">
            <label htmlFor="description">Descripción (Opcional)</label>
            <textarea
              id="description"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              className="form-control"
              placeholder="Detalles sobre el entregable, enlaces, etc."
            />
          </div>
          <button type="submit" className="btn btn-primary">
            <Plus size={18} />
            Añadir Tarea
          </button>
        </form>
      </div>

      <main>
        {error && (
          <div className="glass-card" style={{ padding: '2rem', textAlign: 'center', borderColor: 'var(--danger)' }}>
            <AlertCircle size={40} color="var(--danger)" style={{ marginBottom: '1rem' }} />
            <p style={{ color: 'var(--danger)', fontWeight: '600' }}>{error}</p>
            <button 
              onClick={fetchTasks} 
              className="btn btn-primary" 
              style={{ marginTop: '1rem', background: 'var(--danger)', boxShadow: '0 4px 14px var(--danger-glow)' }}
            >
              Reintentar Conexión
            </button>
          </div>
        )}

        {loading && !error && (
          <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', padding: '4rem 0' }}>
            <Loader2 className="animate-spin" size={32} color="var(--primary)" style={{ animation: 'spin 1s linear infinite' }} />
            <style>{`
              @keyframes spin {
                from { transform: rotate(0deg); }
                to { transform: rotate(360deg); }
              }
            `}</style>
          </div>
        )}

        {!loading && !error && tasks.length === 0 && (
          <div className="glass-card empty-state">
            <ListTodo size={48} className="empty-icon" />
            <h3>No tienes tareas pendientes</h3>
            <p>Comienza agregando una tarea en el formulario de arriba.</p>
          </div>
        )}

        {!loading && !error && tasks.length > 0 && (
          <div className="tasks-grid">
            {tasks.map((task) => (
              <TaskCard
                key={task.id}
                task={task}
                onToggle={handleToggleTask}
                onDelete={handleDeleteTask}
              />
            ))}
          </div>
        )}
      </main>
    </div>
  );
}
