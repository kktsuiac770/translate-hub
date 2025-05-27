import React, { useEffect, useState } from 'react';
import { 
  Container, Typography, Box, Button, TextField,
  Table, TableBody, TableCell, TableContainer, TableHead, 
  TableRow, Paper, Chip 
} from '@mui/material';
import { useParams } from 'react-router-dom';

interface TaskDetailProps {
  user: string;
}

interface Dialogue {
  id: number;
  text: string;
  trans: string;
  translator?: string;
}

interface Change {
  id: number;
  dialogue_id: number;
  user: string;
  new_trans: string;
  status: string;
}

interface Task {
  id: number;
  name: string;
  creator: string;
  status: string;
  dialogues: Dialogue[];
  changes: Change[];
}

const API_BASE = process.env.REACT_APP_API_BASE || 'http://localhost:8080';

const TaskDetail: React.FC<TaskDetailProps> = ({ user }) => {
  const { taskId } = useParams();
  const [task, setTask] = useState<Task | null>(null);
  const [changes, setChanges] = useState<{ [id: number]: string }>({});
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetch(`${API_BASE}/tasks/${taskId}`)
      .then(res => {
        if (res.status === 404) {
          throw new Error(`Task ${taskId} not found`);
        }
        if (!res.ok) {
          throw new Error(`Server error: ${res.status}`);
        }
        return res.json();
      })
      .then(data => {
        if (!data) {
          throw new Error('No data received from server');
        }
        if (!data.id || !data.name) {
          console.warn('Task data missing required fields:', data);
        }
        // Initialize task with empty arrays and default values
        setTask({
          id: data.id || parseInt(taskId as string),
          name: data.name || 'Unnamed Task',
          creator: data.creator || 'Unknown',
          status: data.status || 'Unknown',
          dialogues: Array.isArray(data.dialogues) ? data.dialogues : [],
          changes: Array.isArray(data.changes) ? data.changes : []
        });
        setError(null);
      })
      .catch(err => {
        console.error('Error fetching task:', err);
        setError(err.message);
      });
  }, [taskId]);

  const handleChange = (id: number, value: string) => {
    setChanges({ ...changes, [id]: value });
  };

  const handleSubmit = (id: number) => {
    fetch(`${API_BASE}/changes`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      body: `task_id=${taskId}&dialogue_id=${id}&user=${encodeURIComponent(user)}&new_trans=${encodeURIComponent(changes[id] || '')}`
    }).then(() => {
      setChanges({ ...changes, [id]: '' });
      // Refresh task data to get updated translations
      fetch(`${API_BASE}/tasks/${taskId}`)
        .then(res => res.json())
        .then(data => setTask(data));
    });
  };

  if (error) {
    return (
      <Container>
        <Box mt={4}>
          <Typography color="error">Error: {error}</Typography>
        </Box>
      </Container>
    );
  }

  if (!task) {
    return (
      <Container>
        <Box mt={4}>
          <Typography>Loading...</Typography>
        </Box>
      </Container>
    );
  }

  return (
    <Container>
      <Box mt={4}>
        <Typography variant="h4" gutterBottom>
          Task: {task.name || 'Unnamed Task'} (ID: {taskId})
        </Typography>
        <Typography variant="subtitle1" gutterBottom>
          Created by: {task.creator || 'Unknown'} | Status: {task.status || 'Unknown'}
        </Typography>
        
        <TableContainer component={Paper} sx={{ mt: 3 }}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Original Text</TableCell>
                <TableCell>Translation</TableCell>
                <TableCell>Suggested Changes</TableCell>
                <TableCell>Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {(task.dialogues || []).map(d => (
                <TableRow key={d.id}>
                  <TableCell>{d.text}</TableCell>
                  <TableCell>
                    {d.trans}
                    {d.translator && (
                      <Chip
                        label={`Translated by ${d.translator}`}
                        size="small"
                        sx={{ ml: 1 }}
                      />
                    )}
                  </TableCell>
                  <TableCell>
                    {(task.changes || [])
                      .filter(c => c.dialogue_id === d.id && c.status === 'pending')
                      .map(c => (
                        <Chip
                          key={c.id}
                          label={`${c.user}: ${c.new_trans}`}
                          size="small"
                          sx={{ m: 0.5 }}
                        />
                      ))}
                  </TableCell>
                  <TableCell>
                    <Box display="flex" alignItems="center" gap={1}>
                      <TextField
                        label="Suggest Translation"
                        value={changes[d.id] || ''}
                        onChange={e => handleChange(d.id, e.target.value)}
                        size="small"
                        sx={{ width: 200 }}
                      />
                      <Button 
                        variant="contained" 
                        onClick={() => handleSubmit(d.id)} 
                        disabled={!changes[d.id] || !changes[d.id].trim()}
                        size="small"
                      >
                        Submit
                      </Button>
                    </Box>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      </Box>
    </Container>
  );
};

export default TaskDetail;
