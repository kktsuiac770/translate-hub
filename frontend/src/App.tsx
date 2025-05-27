import React, { useState } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import Login from './Login';
import ProjectList from './ProjectList';
import ProjectDetail from './ProjectDetail';
import TaskDetail from './TaskDetail';
import ReviewChanges from './ReviewChanges';
import './App.css';

function App() {
  const [user, setUser] = useState<string | null>(null);

  if (!user) {
    return <Login onLogin={setUser} />;
  }

  return (
    <Router>
      <Routes>
        <Route path="/" element={<ProjectList user={user} />} />
        <Route path="/project/:projectId" element={<ProjectDetail user={user} />} />
        <Route path="/task/:taskId" element={<TaskDetail user={user} />} />
        <Route path="/task/:taskId/review" element={<ReviewChanges user={user} />} />
        <Route path="*" element={<Navigate to="/" />} />
      </Routes>
    </Router>
  );
}

export default App;