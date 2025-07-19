import React, { useState } from 'react';
import { View, Text, TextInput, TouchableOpacity, StyleSheet, Alert, ActivityIndicator } from 'react-native';
import { useRouter } from 'expo-router';
import { useAuth } from '../contexts/AuthContext';

const API_URL = 'http://192.168.20.76:3000/v1/users/login';

interface LoginResponse {
  access_token: string;
  user?: {
    id: number;
    name: string;
    email: string;
    username?: string;
    avatar?: string;
  };
}

export default function LoginScreen() {
  const router = useRouter();
  const { login } = useAuth();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);

  const handleLogin = async () => {
    if (!email || !password) {
      Alert.alert('Lỗi', 'Vui lòng nhập đầy đủ thông tin');
      return;
    }
    setLoading(true);
    try {
      const res = await fetch(API_URL, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password }),
      });
      const contentType = res.headers.get('content-type');
      let data: LoginResponse;
      if (contentType && contentType.includes('application/json')) {
        data = await res.json();
      } else {
        const textData = await res.text();
        throw new Error(textData || 'Đăng nhập thất bại');
      }
      if (!res.ok) {
        throw new Error((data as any)?.message || 'Đăng nhập thất bại');
      }
      
      // Sử dụng AuthContext để login
      if (data.access_token) {
        const userData = data.user ? {
          ...data.user,
          rating: 4.8,
          memberSince: 2022,
        } : undefined;
        
        await login(data.access_token, userData);
        console.log('TOKEN:', data.access_token);
        Alert.alert('Thành công', 'Đăng nhập thành công!');
        router.replace('/');
      }
    } catch (err) {
      const error = err instanceof Error ? err : new Error(String(err));
      Alert.alert('Lỗi', error.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <View style={styles.container}>
      <Text style={styles.logo}>ShopOnline</Text>
      <Text style={styles.title}>Đăng nhập tài khoản</Text>
      <TextInput
        style={styles.input}
        placeholder="Email"
        autoCapitalize="none"
        keyboardType="email-address"
        value={email}
        onChangeText={setEmail}
      />
      <TextInput
        style={styles.input}
        placeholder="Mật khẩu"
        secureTextEntry
        value={password}
        onChangeText={setPassword}
      />
      <TouchableOpacity style={styles.button} onPress={handleLogin} disabled={loading}>
        {loading ? <ActivityIndicator color="#fff" /> : <Text style={styles.buttonText}>Đăng nhập</Text>}
      </TouchableOpacity>
      <TouchableOpacity onPress={() => router.push('/register')}>
        <Text style={styles.link}>Chưa có tài khoản? Đăng ký</Text>
      </TouchableOpacity>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
    alignItems: 'center',
    justifyContent: 'center',
    padding: 24,
  },
  logo: {
    fontSize: 40,
    fontWeight: 'bold',
    color: '#1976d2',
    marginBottom: 16,
  },
  title: {
    fontSize: 20,
    fontWeight: '600',
    marginBottom: 24,
    color: '#222',
  },
  input: {
    width: '100%',
    height: 48,
    borderColor: '#eee',
    borderWidth: 1,
    borderRadius: 8,
    paddingHorizontal: 16,
    marginBottom: 16,
    backgroundColor: '#fafafa',
  },
  button: {
    width: '100%',
    height: 48,
    backgroundColor: '#1976d2',
    borderRadius: 8,
    alignItems: 'center',
    justifyContent: 'center',
    marginBottom: 16,
  },
  buttonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: 'bold',
  },
  link: {
    color: '#1976d2',
    fontSize: 15,
    marginTop: 8,
  },
});