import React, { useEffect, useState } from 'react';
import { Container, Typography, Box, List, ListItem, ListItemText, Button } from '@mui/material';
import { useParams } from 'react-router-dom';

interface Change {
  id: number;
  dialogue_id: number;
  user: string;
  new_trans: string;
  status: string;
}
const API_BASE = process.env.REACT_APP_API_BASE || 'http://localhost:8080';

const ReviewChanges: React.FC<{ user: string }> = ({ user }) => {
  const { taskId } = useParams();
  const [changes, setChanges] = useState<Change[]>([]);

  useEffect(() => {
    fetch(`${API_BASE}/tasks/${taskId}/changes?user=${encodeURIComponent(user)}`)
      .then(res => res.json())
      .then(setChanges)
      .catch(() => setChanges([]));
  }, [taskId, user]);

  const handleReview = (changeId: number, status: string) => {
    fetch(`${API_BASE}/review`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      body: `task_id=${taskId}&change_id=${changeId}&user=${encodeURIComponent(user)}&status=${status}`
    }).then(() => {
      setChanges(changes.map(c => c.id === changeId ? { ...c, status } : c));
    });
  };

  return (
    <Container>
      <Box mt={4}>
        <Typography variant="h5">Review Changes for Task {taskId}</Typography>
        <List>
          {changes.map(change => (
            <ListItem key={change.id} alignItems="flex-start">
              <ListItemText
                primary={`User: ${change.user}`}
                secondary={`Suggestion: ${change.new_trans} | Status: ${change.status}`}
              />
              {change.status === 'pending' && (
                <>
                  <Button onClick={() => handleReview(change.id, 'approved')} color="success" sx={{ mr: 1 }}>
                    Approve
                  </Button>
                  <Button onClick={() => handleReview(change.id, 'rejected')} color="error">
                    Reject
                  </Button>
                </>
              )}
            </ListItem>
          ))}
        </List>
      </Box>
    </Container>
  );
};

export default ReviewChanges;
