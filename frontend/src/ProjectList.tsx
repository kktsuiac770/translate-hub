import React, { useState, useEffect } from 'react';
import { Container, Typography, Button, Box, List, ListItem, ListItemText, Dialog, DialogTitle, DialogContent, TextField, DialogActions } from '@mui/material';
import { useNavigate } from 'react-router-dom';

interface ProjectListProps {
  user: string;
}

interface Project {
  id: number;
  name: string;
  source_lang?: string;
  target_lang?: string;
}

const API_BASE = process.env.REACT_APP_API_BASE || 'http://localhost:8080';

const ProjectList: React.FC<ProjectListProps> = ({ user }) => {
  const [projects, setProjects] = useState<Project[]>([]);
  const [open, setOpen] = useState(false);
  const [newProjectName, setNewProjectName] = useState('');
  const [newSourceLang, setNewSourceLang] = useState('');
  const [newTargetLang, setNewTargetLang] = useState('');
  const navigate = useNavigate();

  useEffect(() => {
    fetch(`${API_BASE}/projects`)
      .then(res => res.json())
      .then(data => setProjects(Array.isArray(data) ? data : []))
      .catch(() => setProjects([]));
  }, []);

  const handleCreate = () => {
    fetch(`${API_BASE}/projects`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name: newProjectName, source_lang: newSourceLang, target_lang: newTargetLang })
    })
      .then(res => res.json())
      .then(project => {
        setProjects([...projects, project]);
        setOpen(false);
        setNewProjectName('');
        setNewSourceLang('');
        setNewTargetLang('');
      });
  };

  return (
    <Container>
      <Box mt={4}>
        <Typography variant="h4">Projects</Typography>
        <Button variant="contained" color="primary" sx={{ mt: 2 }} onClick={() => setOpen(true)}>
          Create Project
        </Button>
        <List>
          {projects.map(p => (
            <ListItem key={p.id} onClick={() => navigate(`/project/${p.id}`)} style={{ cursor: 'pointer' }}>
              <ListItemText primary={p.name} />
            </ListItem>
          ))}
        </List>
        <Dialog open={open} onClose={() => setOpen(false)}>
          <DialogTitle>New Project</DialogTitle>
          <DialogContent>
            <TextField label="Project Name" value={newProjectName} onChange={e => setNewProjectName(e.target.value)} fullWidth sx={{ mb: 2 }} />
            <TextField label="Source Language (e.g. en)" value={newSourceLang} onChange={e => setNewSourceLang(e.target.value)} fullWidth sx={{ mb: 2 }} />
            <TextField label="Target Language (e.g. jp)" value={newTargetLang} onChange={e => setNewTargetLang(e.target.value)} fullWidth sx={{ mb: 2 }} />
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setOpen(false)}>Cancel</Button>
            <Button onClick={handleCreate} disabled={!newProjectName.trim()}>Create</Button>
          </DialogActions>
        </Dialog>
      </Box>
    </Container>
  );
};

export default ProjectList;
