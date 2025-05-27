import React, { useEffect, useState, useRef } from 'react';
import { Container, Typography, Button, Box, List, ListItem, ListItemText, Dialog, DialogTitle, DialogContent, TextField, DialogActions, Input, Paper } from '@mui/material';
import { useParams, useNavigate } from 'react-router-dom';

interface ProjectDetailProps {
  user: string;
}

interface Task {
  id: number;
  name: string;
  status: string;
}

const API_BASE = process.env.REACT_APP_API_BASE || 'http://localhost:8080';

const ProjectDetail: React.FC<ProjectDetailProps> = ({ user }) => {
  const { projectId } = useParams();
  const [tasks, setTasks] = useState<Task[]>([]);
  const [open, setOpen] = useState(false);
  const [newTaskName, setNewTaskName] = useState('');
  const [file, setFile] = useState<File | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const navigate = useNavigate();

  useEffect(() => {
    fetch(`${API_BASE}/projects/${projectId}/tasks`)
      .then(res => res.json())
      .then(setTasks)
      .catch(() => setTasks([]));
  }, [projectId]);

  const handleCreate = () => {
    if (!file) return;
    const formData = new FormData();
    formData.append('file', file);
    formData.append('name', newTaskName);
    formData.append('creator', user);
    formData.append('project_id', projectId || '');
    fetch(`${API_BASE}/projects/${projectId}/tasks`, {
      method: 'POST',
      body: formData
    })
      .then(res => res.json())
      .then(task => {
        setTasks([...tasks, task]);
        setOpen(false);
        setNewTaskName('');
        setFile(null);
      });
  };

  const handleDrop = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    if (e.dataTransfer.files && e.dataTransfer.files[0]) {
      const droppedFile = e.dataTransfer.files[0];
      setFile(droppedFile);
      setNewTaskName(droppedFile.name.replace(/\.[^/.]+$/, ''));
    }
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const target = e.target as HTMLInputElement;
    if (target.files && target.files[0]) {
      setFile(target.files[0]);
      setNewTaskName(target.files[0].name.replace(/\.[^/.]+$/, ''));
    }
  };

  return (
    <Container>
      <Box mt={4}>
        <Typography variant="h4">Project Detail (ID: {projectId})</Typography>
        <Button variant="contained" color="primary" sx={{ mt: 2 }} onClick={() => setOpen(true)}>
          Create Task
        </Button>
        <List>
          {tasks.map(t => (
            <ListItem key={t.id} onClick={() => navigate(`/task/${t.id}`)} style={{ cursor: 'pointer' }}>
              <ListItemText primary={t.name} secondary={t.status} />
            </ListItem>
          ))}
        </List>
        <Dialog open={open} onClose={() => setOpen(false)}>
          <DialogTitle>New Task</DialogTitle>
          <DialogContent>
            <Paper
              variant="outlined"
              sx={{ p: 2, mb: 2, textAlign: 'center', borderStyle: 'dashed', background: '#fafafa' }}
              onDrop={handleDrop}
              onDragOver={e => e.preventDefault()}
              onClick={() => fileInputRef.current?.click()}
            >
              {file ? (
                <Typography>{file.name}</Typography>
              ) : (
                <Typography color="textSecondary">Drag & drop a .txt file here, or click to select</Typography>
              )}
              <Input
                type="file"
                inputProps={{ accept: '.txt' }}
                onChange={handleFileChange}
                inputRef={fileInputRef}
                sx={{ display: 'none' }}
              />
            </Paper>
            <TextField label="Task Name" value={newTaskName} onChange={e => setNewTaskName(e.target.value)} fullWidth sx={{ mb: 2 }} />
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setOpen(false)}>Cancel</Button>
            <Button onClick={handleCreate} disabled={!newTaskName.trim()}>Create</Button>
          </DialogActions>
        </Dialog>
      </Box>
    </Container>
  );
};

export default ProjectDetail;
