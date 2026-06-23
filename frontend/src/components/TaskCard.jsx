import React from 'react';
import { CheckCircle2, Circle, Trash2, Calendar } from 'lucide-react';

export default function TaskCard({ task, onToggle, onDelete }) {
  const formatDate = (dateString) => {
    if (!dateString) return '';
    const date = new Date(dateString);
    return date.toLocaleDateString('es-ES', {
      day: '2-digit',
      month: 'short',
      year: 'numeric',
    });
  };

  return (
    <div className={`glass-card task-card ${task.completed ? 'completed-card' : ''}`}>
      <div className="task-header">
        <h3 className={`task-title ${task.completed ? 'completed' : ''}`}>
          {task.title}
        </h3>
        <span className={`badge ${task.completed ? 'badge-completed' : 'badge-pending'}`}>
          {task.completed ? 'Completado' : 'Pendiente'}
        </span>
      </div>

      <p className={`task-desc ${task.completed ? 'completed' : ''}`}>
        {task.description || 'Sin descripción adicional.'}
      </p>

      <div className="task-footer">
        <div className="task-date">
          <Calendar size={14} style={{ marginRight: '4px', verticalAlign: 'middle' }} />
          <span>{formatDate(task.created_at)}</span>
        </div>
        <div className="task-actions">
          <button 
            onClick={() => onToggle(task)}
            className="btn-icon" 
            title={task.completed ? "Marcar como pendiente" : "Marcar como completada"}
          >
            {task.completed ? (
              <CheckCircle2 size={18} color="#10b981" />
            ) : (
              <Circle size={18} />
            )}
          </button>
          <button 
            onClick={() => onDelete(task.id)}
            className="btn-icon btn-danger-hover" 
            title="Eliminar tarea"
          >
            <Trash2 size={18} />
          </button>
        </div>
      </div>
    </div>
  );
}
