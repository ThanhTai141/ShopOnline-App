import React, { createContext, useContext, useState, useEffect } from 'react';
import AsyncStorage from '@react-native-async-storage/async-storage';

interface User {
  id: number;
  name: string;
  email: string;
  username?: string;
  avatar?: string;
  rating?: number;
  memberSince?: number;
}

interface AuthContextType {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (token: string, userData?: User) => Promise<void>;
  logout: () => Promise<void>;
  updateUser: (userData: User) => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    loadAuthState();
  }, []);

  const loadAuthState = async () => {
    try {
      const token = await AsyncStorage.getItem('access_token');
      const userData = await AsyncStorage.getItem('user_data');
      
      if (token && userData) {
        try {
          const parsedUser = JSON.parse(userData);
          setUser(parsedUser);
        } catch (error) {
          console.error('Error parsing user data:', error);
          // Nếu user data bị lỗi, thử fetch từ API
          await fetchUserFromAPI(token);
        }
      } else if (token) {
        // Có token nhưng không có user data, thử fetch từ API
        await fetchUserFromAPI(token);
      }
    } catch (error) {
      console.error('Error loading auth state:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const fetchUserFromAPI = async (token: string) => {
    try {
      const response = await fetch('http://192.168.20.76:3000/v1/users/me', {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
      });

      if (response.ok) {
        const userData = await response.json();
        setUser(userData);
        await AsyncStorage.setItem('user_data', JSON.stringify(userData));
      } else {
        console.log('API not available, using default user');
        setDefaultUser();
      }
    } catch (error) {
      console.log('Network error, using default user');
      setDefaultUser();
    }
  };

  const setDefaultUser = () => {
    const defaultUser: User = {
      id: 1,
      name: 'Nguyễn Văn A',
      email: 'user@example.com',
      username: 'nguyenvana123',
      avatar: 'https://randomuser.me/api/portraits/men/32.jpg',
      rating: 4.8,
      memberSince: 2022,
    };
    setUser(defaultUser);
    AsyncStorage.setItem('user_data', JSON.stringify(defaultUser));
  };

  const login = async (token: string, userData?: User) => {
    await AsyncStorage.setItem('access_token', token);
    
    if (userData) {
      setUser(userData);
      await AsyncStorage.setItem('user_data', JSON.stringify(userData));
    } else {
      // Nếu không có userData, tạo mock data từ token hoặc fetch từ API
      await fetchUserFromAPI(token);
    }
  };

  const logout = async () => {
    await AsyncStorage.removeItem('access_token');
    await AsyncStorage.removeItem('user_data');
    setUser(null);
  };

  const updateUser = (userData: User) => {
    setUser(userData);
    AsyncStorage.setItem('user_data', JSON.stringify(userData));
  };

  const value: AuthContextType = {
    user,
    isAuthenticated: !!user,
    isLoading,
    login,
    logout,
    updateUser,
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
}; 