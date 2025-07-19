import React, {  useState } from 'react';
import { View, Text, FlatList, Image, TouchableOpacity, StyleSheet, ActivityIndicator, Dimensions, TextInput, Alert } from 'react-native';
import { useRouter } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { useCart } from '../../contexts/CartContext';
import { useFocusEffect } from '@react-navigation/native';

const API_URL = 'http://192.168.20.76:3000/v1/products';
const numColumns = 2;
const { width } = Dimensions.get('window');
const itemWidth = (width - 36) / numColumns;

export default function HomeScreen() {
  const router = useRouter();
  const { addToCart, cartCount } = useCart();
  const [products, setProducts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [search, setSearch] = useState('');

  useFocusEffect(
    React.useCallback(() => {
      fetchProducts();
    }, [])
  );

  const fetchProducts = async () => {
    setLoading(true);
    setError('');
    try {
      const res = await fetch(API_URL);
      const data = await res.json();
      if (!res.ok) throw new Error(data.message || 'Lỗi tải sản phẩm');
      setProducts(data);
    } catch (err) {
      const error = err instanceof Error ? err : new Error(String(err));
      setError(error.message);
    } finally {
      setLoading(false);
    }
  };

  const handleAddToCart = (product) => {
    addToCart({
      id: product.id,
      name: product.name,
      price: product.price,
      image_url: product.image_url,
      quantity: 1,
    });
    Alert.alert('Thành công', 'Đã thêm sản phẩm vào giỏ hàng!');
  };

  const filteredProducts = products.filter(p =>
    p.name.toLowerCase().includes(search.toLowerCase())
  );

  const renderItem = ({ item }) => (
    <TouchableOpacity
      style={styles.item}
      activeOpacity={0.8}
      onPress={() => router.push(`/product/${item.id}`)}
    >
      <Image source={{ uri: item.image_url }} style={styles.image} resizeMode="cover" />
      <Text style={styles.name} numberOfLines={2}>{item.name}</Text>
      <Text style={styles.category}>{item.category}</Text>
      <Text style={styles.desc} numberOfLines={1}>{item.description}</Text>
      <Text style={styles.price}>{item.price.toLocaleString()} đ</Text>
      <TouchableOpacity 
        style={styles.cartBtn} 
        activeOpacity={0.7}
        onPress={() => handleAddToCart(item)}
      >
        <Ionicons name="cart-outline" size={20} color="#1976d2" />
      </TouchableOpacity>
    </TouchableOpacity>
  );

  if (loading) {
    return <View style={styles.center}><ActivityIndicator size="large" color="#1976d2" /></View>;
  }
  if (error) {
    return <View style={styles.center}><Text style={{ color: 'red' }}>{error}</Text></View>;
  }

  return (
    <View style={styles.container}>
      <View style={styles.header}>
        <View style={styles.searchWrap}>
          <Ionicons name="search" size={20} color="#888" style={{ marginLeft: 12 }} />
          <TextInput
            style={styles.search}
            placeholder="Tìm kiếm sản phẩm..."
            value={search}
            onChangeText={setSearch}
            placeholderTextColor="#aaa"
          />
        </View>
        <TouchableOpacity 
          style={styles.cartButton} 
          activeOpacity={0.7}
          onPress={() => router.push('/cart')}
        >
          <Ionicons name="cart" size={24} color="#1976d2" />
          {cartCount > 0 && (
            <View style={styles.cartBadge}>
              <Text style={styles.cartBadgeText}>{cartCount > 99 ? '99+' : cartCount}</Text>
            </View>
          )}
        </TouchableOpacity>
      </View>
      <FlatList
        data={filteredProducts}
        renderItem={renderItem}
        keyExtractor={item => item.id.toString()}
        numColumns={numColumns}
        contentContainerStyle={{ paddingBottom: 16 }}
        columnWrapperStyle={{ gap: 12, marginBottom: 16 }}
        showsVerticalScrollIndicator={false}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
    paddingHorizontal: 12,
    paddingTop: 16,
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 18,
    gap: 12,
  },
  searchWrap: {
    flex: 1,
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: '#f5f5f5',
    borderRadius: 24,
    paddingVertical: 6,
    paddingRight: 12,
    shadowColor: '#000',
    shadowOpacity: 0.04,
    shadowRadius: 2,
    elevation: 1,
  },
  search: {
    flex: 1,
    height: 38,
    fontSize: 16,
    marginLeft: 8,
    color: '#222',
  },
  cartButton: {
    width: 44,
    height: 44,
    borderRadius: 22,
    backgroundColor: '#f5f5f5',
    alignItems: 'center',
    justifyContent: 'center',
    position: 'relative',
    shadowColor: '#000',
    shadowOpacity: 0.04,
    shadowRadius: 2,
    elevation: 1,
  },
  cartBadge: {
    position: 'absolute',
    top: -2,
    right: -2,
    backgroundColor: '#e53935',
    borderRadius: 10,
    minWidth: 20,
    height: 20,
    alignItems: 'center',
    justifyContent: 'center',
    paddingHorizontal: 4,
  },
  cartBadgeText: {
    color: '#fff',
    fontSize: 12,
    fontWeight: 'bold',
  },
  item: {
    backgroundColor: '#fff',
    borderRadius: 14,
    padding: 10,
    alignItems: 'center',
    width: itemWidth,
    shadowColor: '#000',
    shadowOpacity: 0.08,
    shadowRadius: 8,
    elevation: 2,
    position: 'relative',
  },
  image: {
    width: itemWidth - 20,
    height: itemWidth - 20,
    borderRadius: 10,
    marginBottom: 8,
    backgroundColor: '#eee',
  },
  name: {
    fontSize: 15,
    fontWeight: '500',
    color: '#222',
    marginBottom: 4,
    textAlign: 'center',
  },
  price: {
    fontSize: 16,
    fontWeight: 'bold',
    color: '#1976d2',
    marginBottom: 6,
  },
  cartBtn: {
    position: 'absolute',
    bottom: 10,
    right: 10,
    backgroundColor: '#fff',
    borderRadius: 16,
    padding: 6,
    shadowColor: '#000',
    shadowOpacity: 0.08,
    shadowRadius: 4,
    elevation: 2,
  },
  center: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
  },
  category: {
    fontSize: 12,
    color: '#888',
    marginBottom: 2,
    textAlign: 'center',
  },
  desc: {
    fontSize: 13,
    color: '#444',
    marginBottom: 4,
    textAlign: 'center',
  },
});