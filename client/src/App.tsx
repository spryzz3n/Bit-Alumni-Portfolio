import { BrowserRouter, Route, Routes } from 'react-router-dom'
import StatusPage from './features/status/StatusPage'

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<StatusPage />} />
      </Routes>
    </BrowserRouter>
  )
}
