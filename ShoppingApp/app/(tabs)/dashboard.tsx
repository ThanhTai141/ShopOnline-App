import React, { useEffect, useState } from 'react';
import { View, Text, FlatList, TouchableOpacity, StyleSheet, Alert, TextInput, Modal, Button } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import AsyncStorage from '@react-native-async-storage/async-storage';

const API_URL = 'http://192.168.20.76:3000/v1/products';

interface Product {
  id: number;
  name: string;
  price: number;
  image_url: string;
  category: string;
  description: string;
}

export default function DashboardScreen() {
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [modalVisible, setModalVisible] = useState(false);
  const [editProduct, setEditProduct] = useState<Product | null>(null);
  const [form, setForm] = useState<Omit<Product, 'id'>>({ name: '', price: 0, image_url: '', category: '', description: '' });

  useEffect(() => {
    fetchProducts();
  }, []);

  const fetchProducts = async () => {
    setLoading(true);
    try {
      const res = await fetch(API_URL);
      const data = await res.json();
      setProducts(data);
    } catch (err) {
      Alert.alert('Lỗi', 'Không thể tải sản phẩm');
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id: number) => {
    Alert.alert('Xác nhận', 'Bạn có chắc muốn xoá sản phẩm này?', [
      { text: 'Huỷ', style: 'cancel' },
      {
        text: 'Xoá', style: 'destructive', onPress: async () => {
          const token = await AsyncStorage.getItem('access_token');
          const res = await fetch(`${API_URL}/${id}`, {
            method: 'DELETE',
            headers: {
              'Authorization': 'Bearer ' + token,
            },
          });
          if (!res.ok) {
            let errMsg = 'Có lỗi xảy ra';
            try { errMsg = (await res.json()).error; } catch {}
            Alert.alert('Lỗi', errMsg);
            return;
          }
          fetchProducts();
        }
      }
    ]);
  };

  const handleEdit = (product: Product) => {
    setEditProduct(product);
    setForm({
      name: product.name,
      price: product.price,
      image_url: product.image_url,
      category: product.category,
      description: product.description,
    });
    setModalVisible(true);
  };

  const handleAdd = () => {
    setEditProduct(null);
    setForm({ name: '', price: 0, image_url: '', category: '', description: '' });
    setModalVisible(true);
  };

  const handleSubmit = async () => {
    if (!form.name || !form.price) {
      Alert.alert('Lỗi', 'Vui lòng nhập đầy đủ tên và giá');
      return;
    }
    const token = await AsyncStorage.getItem('access_token');
    if (editProduct) {
      // Update
      const res = await fetch(`${API_URL}/${editProduct.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer ' + token,
        },
        body: JSON.stringify(form),
      });
      if (!res.ok) {
        let errMsg = 'Có lỗi xảy ra';
        try { errMsg = (await res.json()).error; } catch {}
        Alert.alert('Lỗi', errMsg);
        return;
      }
    } else {
      // Create
      const res = await fetch(API_URL, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer ' + token,
        },
        body: JSON.stringify(form),
      });
      if (!res.ok) {
        let errMsg = 'Có lỗi xảy ra';
        try { errMsg = (await res.json()).error; } catch {}
        Alert.alert('Lỗi', errMsg);
        return;
      }
    }
    setModalVisible(false);
    fetchProducts();
  };

  return (
    <View style={styles.container}>
      <Text style={styles.title}>Quản lý sản phẩm</Text>
      <TouchableOpacity style={styles.addBtn} onPress={handleAdd}>
        <Ionicons name="add-circle" size={24} color="#1976d2" />
        <Text style={styles.addBtnText}>Thêm sản phẩm</Text>
      </TouchableOpacity>
      <FlatList
        data={products}
        keyExtractor={item => item.id.toString()}
        refreshing={loading}
        onRefresh={fetchProducts}
        renderItem={({ item }: { item: Product }) => (
          <View style={styles.itemRow}>
            <View style={{ flex: 1 }}>
              <Text style={styles.itemName}>{item.name}</Text>
              <Text style={styles.itemPrice}>{item.price.toLocaleString()} đ</Text>
            </View>
            <TouchableOpacity onPress={() => handleEdit(item)} style={styles.iconBtn}>
              <Ionicons name="create-outline" size={22} color="#1976d2" />
            </TouchableOpacity>
            <TouchableOpacity onPress={() => handleDelete(item.id)} style={styles.iconBtn}>
              <Ionicons name="trash-outline" size={22} color="#e53935" />
            </TouchableOpacity>
          </View>
        )}
      />

      {/* Modal Thêm/Sửa */}
      <Modal visible={modalVisible} animationType="slide" transparent>
        <View style={styles.modalOverlay}>
          <View style={styles.modalContent}>
            <Text style={styles.modalTitle}>{editProduct ? 'Sửa sản phẩm' : 'Thêm sản phẩm'}</Text>
            <TextInput style={styles.input} placeholder="Tên sản phẩm" value={form.name} onChangeText={v => setForm(f => ({ ...f, name: v }))} />
            <TextInput style={styles.input} placeholder="Giá" value={form.price ? form.price.toString() : ''} keyboardType="numeric" onChangeText={v => setForm(f => ({ ...f, price: Number(v) }))} />
            <TextInput style={styles.input} placeholder="Image URL" value={form.image_url} onChangeText={v => setForm(f => ({ ...f, image_url: v }))} />
            <TextInput style={styles.input} placeholder="Danh mục" value={form.category} onChangeText={v => setForm(f => ({ ...f, category: v }))} />
            <TextInput style={styles.input} placeholder="Mô tả" value={form.description} onChangeText={v => setForm(f => ({ ...f, description: v }))} />
            <View style={{ flexDirection: 'row', justifyContent: 'flex-end', marginTop: 12 }}>
              <Button title="Huỷ" color="#888" onPress={() => setModalVisible(false)} />
              <View style={{ width: 12 }} />
              <Button title={editProduct ? 'Cập nhật' : 'Thêm'} color="#1976d2" onPress={handleSubmit} />
            </View>
          </View>
        </View>
      </Modal>
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: '#f6f8fa', padding: 16 },
  title: { fontSize: 22, fontWeight: 'bold', color: '#1976d2', marginBottom: 16 },
  addBtn: { flexDirection: 'row', alignItems: 'center', marginBottom: 16 },
  addBtnText: { color: '#1976d2', fontSize: 16, fontWeight: 'bold', marginLeft: 6 },
  itemRow: { flexDirection: 'row', alignItems: 'center', backgroundColor: '#fff', borderRadius: 10, padding: 12, marginBottom: 10, shadowColor: '#000', shadowOpacity: 0.04, shadowRadius: 4, elevation: 1 },
  itemName: { fontSize: 16, fontWeight: 'bold', color: '#222' },
  itemPrice: { fontSize: 15, color: '#1976d2', marginTop: 2 },
  iconBtn: { padding: 8 },
  modalOverlay: { flex: 1, backgroundColor: 'rgba(0,0,0,0.2)', justifyContent: 'center', alignItems: 'center' },
  modalContent: { backgroundColor: '#fff', borderRadius: 12, padding: 20, width: '90%' },
  modalTitle: { fontSize: 18, fontWeight: 'bold', color: '#1976d2', marginBottom: 12 },
  input: { borderWidth: 1, borderColor: '#eee', borderRadius: 8, padding: 10, marginBottom: 10, backgroundColor: '#fafafa' },
});
