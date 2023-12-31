import { useEffect, useState } from 'react';
import './App.css';

function App() {
  const [greeting, setGreeting] = useState<string | undefined>(undefined);

  useEffect(() => {
    const fetchGreeting = async () => {
      await fetch("/api/greeting/George")
      .then(res => {
        if (!res.ok) {
          return console.log('Network response was not ok');
        }
        return res.json();
      })
      .then(data => {
        setGreeting(data.message);
      });
    };

    fetchGreeting();
  }, [greeting]);

  return (
    <>
      { greeting }
    </>
  );
}

export default App;
