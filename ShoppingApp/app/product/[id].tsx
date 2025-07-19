import React, { useEffect, useState } from 'react';
import { View, Text, Image, StyleSheet, ActivityIndicator, ScrollView, TouchableOpacity, Alert, Dimensions } from 'react-native';
import { useLocalSearchParams } from 'expo-router';
import { Ionicons, MaterialIcons } from '@expo/vector-icons';
import Swiper from 'react-native-swiper';
import { useCart } from '../../contexts/CartContext';

const API_URL = 'http://192.168.20.76:3000/v1/products';
const { width } = Dimensions.get('window');

interface Product {
  id: number;
  name: string;
  price: number;
  image_url: string;
  category: string;
  description: string;
}

const MOCK_IMAGES = [
  'https://images.unsplash.com/photo-1517336714731-489689fd1ca8?auto=format&fit=crop&w=400&q=80',
  'https://images.unsplash.com/photo-1512436991641-6745cdb1723f?auto=format&fit=crop&w=400&q=80',
  'https://images.unsplash.com/photo-1517263904808-5dc0d6a2b43d?auto=format&fit=crop&w=400&q=80',
  'https://images.unsplash.com/photo-1515378791036-0648a3ef77b2?auto=format&fit=crop&w=400&q=80',
];

export default function ProductDetail() {
  const { id } = useLocalSearchParams();
  const { addToCart } = useCart();
  const [product, setProduct] = useState<Product | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [selectedColor, setSelectedColor] = useState('Đen');
  const [quantity, setQuantity] = useState(1);
  const [tab, setTab] = useState('desc');

  useEffect(() => {
    fetchProduct();
  }, [id]);

  const fetchProduct = async () => {
    setLoading(true);
    setError('');
    try {
      const res = await fetch(`${API_URL}/${id}`);
      const data = await res.json();
      if (!res.ok) throw new Error(data.message || 'Không tìm thấy sản phẩm');
      setProduct(data);
    } catch (err) {
      const error = err instanceof Error ? err : new Error(String(err));
      setError(error.message);
    } finally {
      setLoading(false);
    }
  };

  const handleAddToCart = () => {
    if (!product) return;
    
    addToCart({
      id: product.id,
      name: product.name,
      price: product.price,
      image_url: product.image_url,
      quantity: quantity,
      color: selectedColor,
    });
    
    Alert.alert(
      'Thành công', 
      `Đã thêm ${quantity} sản phẩm "${product.name}" (${selectedColor}) vào giỏ hàng!`
    );
  };

  if (loading) {
    return <View style={styles.center}><ActivityIndicator size="large" color="#1976d2" /></View>;
  }
  if (error || !product) {
    return <View style={styles.center}><Text style={{ color: 'red' }}>{error || 'Không tìm thấy sản phẩm'}</Text></View>;
  }

  // Mock data cho các trường chưa có
  const images = [product.image_url, ...MOCK_IMAGES.slice(1)];
  const oldPrice = product.price * 1.25;
  const discount = Math.round(100 - (product.price / oldPrice) * 100);
  const rating = 4.8;
  const ratingCount = 2100;
  const sold = 5200;
  const colors = ['Đen', 'Trắng', 'Xanh dương'];
  const stock = 15;

  return (
    <ScrollView style={styles.container} contentContainerStyle={{ paddingBottom: 32 }}>
      {/* Ảnh sản phẩm dạng slider */}
      <View style={styles.imageWrap}>
        <Swiper style={styles.swiper} showsButtons={false} dotColor="#eee" activeDotColor="#1976d2" height={width-40}>
          {images.map((img, idx) => (
            <Image key={idx} source={{ uri: img }} style={styles.image} resizeMode="cover" />
          ))}
        </Swiper>
      </View>
      {/* Thông tin sản phẩm */}
      <View style={styles.infoCard}>
        <Text style={styles.name}>{product.name}</Text>
        <View style={styles.row}>
          <Ionicons name="star" size={16} color="#f6b100" style={{ marginRight: 2 }} />
          <Text style={styles.rating}>{rating} ({(ratingCount/1000).toFixed(1)}k đánh giá)</Text>
          <Text style={styles.dot}>·</Text>
          <Text style={styles.sold}>Đã bán {Math.round(sold/1000)}k</Text>
        </View>
        <View style={styles.priceRow}>
          <Text style={styles.price}>{product.price.toLocaleString()} đ</Text>
          <Text style={styles.oldPrice}>{oldPrice.toLocaleString()} đ</Text>
          <Text style={styles.discount}>-{discount}%</Text>
        </View>
        <Text style={styles.category}>{product.category}</Text>
      </View>
      {/* Phân loại & số lượng */}
      <View style={styles.optionCard}>
        <Text style={styles.optionLabel}>Phân loại</Text>
        <View style={styles.colorRow}>
          {colors.map(color => (
            <TouchableOpacity
              key={color}
              style={[styles.colorBtn, selectedColor === color && styles.colorBtnActive]}
              onPress={() => setSelectedColor(color)}
            >
              <Text style={[styles.colorBtnText, selectedColor === color && { color: '#1976d2', fontWeight: 'bold' }]}>{color}</Text>
            </TouchableOpacity>
          ))}
        </View>
        <Text style={styles.optionLabel}>Số lượng</Text>
        <View style={styles.qtyRow}>
          <TouchableOpacity style={styles.qtyBtn} onPress={() => setQuantity(q => Math.max(1, q-1))}><Text style={styles.qtyBtnText}>-</Text></TouchableOpacity>
          <Text style={styles.qtyNum}>{quantity}</Text>
          <TouchableOpacity style={styles.qtyBtn} onPress={() => setQuantity(q => Math.min(stock, q+1))}><Text style={styles.qtyBtnText}>+</Text></TouchableOpacity>
          <Text style={styles.qtyStock}>Còn lại {stock} sản phẩm</Text>
        </View>
      </View>
      {/* Thông tin giao hàng */}
      <View style={styles.shipRow}>
        <View style={styles.shipItem}><MaterialIcons name="local-shipping" size={20} color="#1976d2" /><Text style={styles.shipText}>Giao hàng 2-3 ngày</Text></View>
        <View style={styles.shipItem}><Ionicons name="pricetag" size={20} color="#1976d2" /><Text style={styles.shipText}>Freeship toàn quốc</Text></View>
        <View style={styles.shipItem}><Ionicons name="refresh" size={20} color="#1976d2" /><Text style={styles.shipText}>Đổi trả 30 ngày</Text></View>
      </View>
      {/* Tab mô tả/đánh giá */}
      <View style={styles.tabRow}>
        <TouchableOpacity style={[styles.tabBtn, tab==='desc' && styles.tabActive]} onPress={()=>setTab('desc')}><Text style={[styles.tabText, tab==='desc' && { color: '#1976d2', fontWeight: 'bold' }]}>Mô tả</Text></TouchableOpacity>
        <TouchableOpacity style={[styles.tabBtn, tab==='review' && styles.tabActive]} onPress={()=>setTab('review')}><Text style={[styles.tabText, tab==='review' && { color: '#1976d2', fontWeight: 'bold' }]}>Đánh giá ({ratingCount})</Text></TouchableOpacity>
      </View>
      {tab==='desc' ? (
        <View style={styles.tabContent}>
          <Text style={styles.desc}>{product.description}</Text>
          <Text style={styles.desc}>- Chống ồn chủ động ANC{"\n"}- Bluetooth 5.0{"\n"}- Pin dùng 30 giờ</Text>
        </View>
      ) : (
        <View style={styles.tabContent}><Text style={{ color: '#888' }}>Chức năng đánh giá sẽ cập nhật sau.</Text></View>
      )}
      {/* Nút hành động */}
      <View style={styles.actionRow}>
        <TouchableOpacity style={styles.buyBtn} activeOpacity={0.85}><Text style={styles.buyBtnText}>Mua ngay</Text></TouchableOpacity>
        <TouchableOpacity style={styles.addBtnCart} onPress={handleAddToCart} activeOpacity={0.85}>
          <Ionicons name="cart" size={22} color="#1976d2" style={{ marginRight: 6 }} />
          <Text style={styles.addBtnCartText}>Thêm vào giỏ</Text>
        </TouchableOpacity>
      </View>
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: '#fff' },
  imageWrap: { alignItems: 'center', paddingTop: 18, paddingBottom: 8, backgroundColor: '#fff' },
  swiper: { borderRadius: 18, overflow: 'hidden' },
  image: { width: width - 40, height: width - 40, borderRadius: 18, backgroundColor: '#eee', borderWidth: 4, borderColor: '#fff', shadowColor: '#000', shadowOpacity: 0.10, shadowRadius: 16, elevation: 6 },
  infoCard: { backgroundColor: '#fff', marginHorizontal: 16, marginTop: 18, borderRadius: 18, padding: 20, shadowColor: '#000', shadowOpacity: 0.06, shadowRadius: 12, elevation: 2 },
  name: { fontSize: 22, fontWeight: 'bold', color: '#222', marginBottom: 2, textAlign: 'left' },
  row: { flexDirection: 'row', alignItems: 'center', marginBottom: 8 },
  rating: { fontSize: 13, color: '#f6b100', marginRight: 8 },
  dot: { color: '#bbb', marginHorizontal: 4 },
  sold: { fontSize: 13, color: '#888' },
  priceRow: { flexDirection: 'row', alignItems: 'center', marginBottom: 8 },
  price: { fontSize: 22, fontWeight: 'bold', color: '#1976d2', marginRight: 12 },
  oldPrice: { fontSize: 15, color: '#bbb', textDecorationLine: 'line-through', marginRight: 8 },
  discount: { fontSize: 15, color: '#fff', backgroundColor: '#1976d2', borderRadius: 6, paddingHorizontal: 8, paddingVertical: 2, overflow: 'hidden' },
  category: { fontSize: 13, color: '#888', marginBottom: 8, textAlign: 'left' },
  optionCard: { backgroundColor: '#fff', marginHorizontal: 16, marginTop: 18, borderRadius: 18, padding: 20, shadowColor: '#000', shadowOpacity: 0.04, shadowRadius: 8, elevation: 1 },
  optionLabel: { fontSize: 15, fontWeight: 'bold', color: '#222', marginBottom: 8 },
  colorRow: { flexDirection: 'row', gap: 8, marginBottom: 12 },
  colorBtn: { borderWidth: 1, borderColor: '#eee', borderRadius: 8, paddingHorizontal: 16, paddingVertical: 6, marginRight: 8, backgroundColor: '#fafafa' },
  colorBtnActive: { borderColor: '#1976d2', backgroundColor: '#fff' },
  colorBtnText: { fontSize: 14, color: '#444' },
  qtyRow: { flexDirection: 'row', alignItems: 'center', marginBottom: 8 },
  qtyBtn: { width: 32, height: 32, borderRadius: 8, backgroundColor: '#f5f5f5', alignItems: 'center', justifyContent: 'center', marginHorizontal: 4 },
  qtyBtnText: { fontSize: 18, color: '#222' },
  qtyNum: { fontSize: 16, fontWeight: 'bold', color: '#222', marginHorizontal: 8 },
  qtyStock: { fontSize: 13, color: '#888', marginLeft: 12 },
  shipRow: { flexDirection: 'row', justifyContent: 'space-between', marginHorizontal: 16, marginTop: 18, marginBottom: 8 },
  shipItem: { flexDirection: 'row', alignItems: 'center', gap: 4 },
  shipText: { fontSize: 13, color: '#444', marginLeft: 4 },
  tabRow: { flexDirection: 'row', borderBottomWidth: 1, borderColor: '#eee', marginHorizontal: 16, marginTop: 18 },
  tabBtn: { flex: 1, alignItems: 'center', paddingVertical: 12 },
  tabActive: { borderBottomWidth: 2, borderColor: '#1976d2' },
  tabText: { fontSize: 15, color: '#444' },
  tabContent: { paddingHorizontal: 20, paddingVertical: 16 },
  desc: { fontSize: 15, color: '#444', marginBottom: 4, textAlign: 'left', lineHeight: 22 },
  actionRow: { flexDirection: 'row', gap: 12, marginHorizontal: 16, marginTop: 24 },
  buyBtn: { flex: 1, backgroundColor: '#1976d2', borderRadius: 16, alignItems: 'center', justifyContent: 'center', paddingVertical: 16 },
  buyBtnText: { color: '#fff', fontSize: 16, fontWeight: 'bold' },
  addBtnCart: { flex: 1, backgroundColor: '#fff', borderWidth: 2, borderColor: '#1976d2', borderRadius: 16, alignItems: 'center', justifyContent: 'center', flexDirection: 'row', paddingVertical: 16 },
  addBtnCartText: { color: '#1976d2', fontSize: 16, fontWeight: 'bold' },
  center: { flex: 1, alignItems: 'center', justifyContent: 'center' },
});