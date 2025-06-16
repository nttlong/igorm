import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Login from './pages/Login';

const Dashboard = () => {
  return (
    <Box
      sx={{
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        minHeight: '100vh',
        bgcolor: 'grey.100',
      }}
    >
      <Typography variant="h1" color="text.primary">
        Chào mừng đến với Dashboard!
      </Typography>
    </Box>
  );
};

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Login />} />
        <Route path="/dashboard" element={<Dashboard />} />
      </Routes>
    </Router>
  );
}

export default App;