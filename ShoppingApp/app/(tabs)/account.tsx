import React from 'react';
import { View, Text, StyleSheet, TouchableOpacity, Image, ScrollView, Dimensions, ActivityIndicator } from 'react-native';
import { Ionicons, MaterialIcons, FontAwesome5 } from '@expo/vector-icons';
import { useRouter } from 'expo-router';
import { useAuth } from '../../contexts/AuthContext';

const { width } = Dimensions.get('window');

const quickActions = [
  { icon: <MaterialIcons name="local-shipping" size={28} color="#1976d2" />, label: 'Theo dõi đơn hàng' },
  { icon: <Ionicons name="chatbubble-ellipses" size={28} color="#1976d2" />, label: 'Chat với shop' },
  { icon: <MaterialIcons name="verified-user" size={28} color="#1976d2" />, label: 'Bảo hiểm' },
  { icon: <Ionicons name="gift" size={28} color="#1976d2" />, label: 'Khuyến mãi' },
];

const menu = [
  { icon: <Ionicons name="cube" size={22} color="#1976d2" />, label: 'Đơn hàng', badge: 3 },
  { icon: <Ionicons name="heart" size={22} color="#1976d2" />, label: 'Yêu thích', badge: 12 },
  { icon: <Ionicons name="star" size={22} color="#1976d2" />, label: 'Đánh giá' },
  { icon: <Ionicons name="pricetag" size={22} color="#1976d2" />, label: 'Khuyến mãi', badge: 5 },
  { icon: <FontAwesome5 name="wallet" size={22} color="#1976d2" />, label: 'Ví điện tử' },
  { icon: <Ionicons name="location" size={22} color="#1976d2" />, label: 'Địa chỉ' },
  { icon: <Ionicons name="notifications" size={22} color="#1976d2" />, label: 'Thông báo', badge: 8 },
  { icon: <Ionicons name="settings" size={22} color="#1976d2" />, label: 'Cài đặt' },
  { icon: <Ionicons name="help-circle" size={22} color="#1976d2" />, label: 'Hỗ trợ' },
];

export default function AccountScreen() {
  const router = useRouter();
  const { user, isAuthenticated, isLoading, logout } = useAuth();

  const handleLogout = async () => {
    await logout();
    router.replace('/login');
  };

  if (isLoading) {
    return <View style={styles.center}><ActivityIndicator size="large" color="#1976d2" /></View>;
  }

  if (!isAuthenticated) {
    return (
      <View style={{ flex: 1, alignItems: 'center', justifyContent: 'center', backgroundColor: '#f6f8fa' }}>
        <Ionicons name="person-circle-outline" size={80} color="#1976d2" style={{ marginBottom: 16 }} />
        <Text style={{ fontSize: 18, color: '#1976d2', marginBottom: 12 }}>Bạn chưa đăng nhập</Text>
        <TouchableOpacity style={{ backgroundColor: '#1976d2', borderRadius: 8, paddingHorizontal: 32, paddingVertical: 12 }} onPress={() => router.replace('/login')}>
          <Text style={{ color: '#fff', fontSize: 16, fontWeight: 'bold' }}>Đăng nhập ngay</Text>
        </TouchableOpacity>
      </View>
    );
  }

  return (
    <View style={{ flex: 1, backgroundColor: '#f6f8fa' }}>
      <ScrollView contentContainerStyle={{ paddingBottom: 32 }}>
        {/* Header */}
        <View style={styles.headerWrap}>
          <View style={styles.headerRow}>
            <Image 
              source={{ uri: user?.avatar || 'https://randomuser.me/api/portraits/men/32.jpg' }} 
              style={styles.avatar} 
            />
            <View style={{ flex: 1, marginLeft: 16 }}>
              <View style={{ flexDirection: 'row', alignItems: 'center' }}>
                <Text style={styles.name}>{user?.name || 'Đang tải...'}</Text>
                <Ionicons name="pencil" size={16} color="#fff" style={{ marginLeft: 6 }} />
              </View>
              <Text style={styles.username}>@{user?.username || user?.email?.split('@')[0] || 'user'}</Text>
              <View style={{ flexDirection: 'row', alignItems: 'center', marginTop: 4 }}>
                <Ionicons name="star" size={14} color="#ffd600" />
                <Text style={styles.rating}>{user?.rating || 4.8}</Text>
                <Text style={styles.member}>Thành viên từ {user?.memberSince || 2022}</Text>
              </View>
            </View>
          </View>
          <View style={styles.headerStats}>
            <View style={styles.statItem}>
              <Text style={styles.statNum}>158</Text>
              <Text style={styles.statLabel}>Đơn hàng</Text>
            </View>
            <View style={styles.statItem}>
              <Text style={styles.statNum}>2.4K</Text>
              <Text style={styles.statLabel}>Điểm thưởng</Text>
            </View>
            <View style={styles.statItem}>
              <Text style={styles.statNum}>12</Text>
              <Text style={styles.statLabel}>Ưu đãi</Text>
            </View>
          </View>
        </View>
        {/* Quick actions */}
        <View style={styles.quickRow}>
          {quickActions.map((item, idx) => (
            <View key={idx} style={styles.quickItem}>
              <View style={styles.quickIcon}>{item.icon}</View>
              <Text style={styles.quickLabel}>{item.label}</Text>
            </View>
          ))}
        </View>
        {/* Menu */}
        <View style={styles.menuCard}>
          {menu.map((item, idx) => (
            <TouchableOpacity key={item.label} style={styles.menuRow} activeOpacity={0.7}>
              <View style={styles.menuLeft}>
                {item.icon}
                <Text style={styles.menuLabel}>{item.label}</Text>
              </View>
              <View style={styles.menuRight}>
                {item.badge && <View style={styles.badge}><Text style={styles.badgeText}>{item.badge}</Text></View>}
                <Ionicons name="chevron-forward" size={20} color="#bbb" />
              </View>
            </TouchableOpacity>
          ))}
        </View>
        {/* Đăng xuất */}
        <TouchableOpacity style={styles.logoutBtn} onPress={handleLogout} activeOpacity={0.8}>
          <Ionicons name="log-out-outline" size={20} color="#1976d2" style={{ marginRight: 8 }} />
          <Text style={styles.logoutText}>Đăng xuất</Text>
        </TouchableOpacity>
      </ScrollView>
      {/* Nút chat nổi */}
      <TouchableOpacity style={styles.fab} activeOpacity={0.8}>
        <Ionicons name="chatbubble-ellipses" size={28} color="#fff" />
      </TouchableOpacity>
    </View>
  );
}


const styles = StyleSheet.create({
  headerWrap: {
    backgroundColor: '#1976d2',
    borderBottomLeftRadius: 18,
    borderBottomRightRadius: 18,
    paddingTop: 36,
    paddingBottom: 18,
    paddingHorizontal: 18,
  },
  headerRow: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 18,
  },
  avatar: {
    width: 64,
    height: 64,
    borderRadius: 32,
    borderWidth: 2,
    borderColor: '#fff',
    backgroundColor: '#eee',
  },
  name: {
    fontSize: 20,
    fontWeight: 'bold',
    color: '#fff',
  },
  username: {
    fontSize: 14,
    color: '#e3f2fd',
    marginTop: 2,
  },
  rating: {
    fontSize: 13,
    color: '#ffd600',
    marginLeft: 4,
    marginRight: 8,
    fontWeight: 'bold',
  },
  member: {
    fontSize: 13,
    color: '#e3f2fd',
    marginLeft: 8,
  },
  headerStats: {
    flexDirection: 'row',
    backgroundColor: '#fff',
    borderRadius: 12,
    padding: 12,
    justifyContent: 'space-between',
    alignItems: 'center',
    shadowColor: '#000',
    shadowOpacity: 0.04,
    shadowRadius: 4,
    elevation: 1,
  },
  statItem: {
    alignItems: 'center',
    flex: 1,
  },
  statNum: {
    fontSize: 18,
    fontWeight: 'bold',
    color: '#1976d2',
  },
  statLabel: {
    fontSize: 13,
    color: '#888',
    marginTop: 2,
  },
  quickRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    backgroundColor: '#fff',
    borderRadius: 12,
    marginHorizontal: 16,
    marginTop: 18,
    padding: 16,
    shadowColor: '#000',
    shadowOpacity: 0.04,
    shadowRadius: 4,
    elevation: 1,
  },
  quickItem: {
    alignItems: 'center',
    flex: 1,
  },
  quickIcon: {
    backgroundColor: '#e3f2fd',
    borderRadius: 24,
    padding: 10,
    marginBottom: 6,
  },
  quickLabel: {
    fontSize: 13,
    color: '#1976d2',
    textAlign: 'center',
  },
  menuCard: {
    backgroundColor: '#fff',
    borderRadius: 12,
    marginHorizontal: 16,
    marginTop: 18,
    paddingVertical: 4,
    shadowColor: '#000',
    shadowOpacity: 0.04,
    shadowRadius: 4,
    elevation: 1,
  },
  menuRow: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingVertical: 14,
    paddingHorizontal: 16,
    borderBottomWidth: 1,
    borderColor: '#f0f0f0',
  },
  menuLeft: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  menuLabel: {
    fontSize: 15,
    color: '#222',
    marginLeft: 12,
  },
  menuRight: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  badge: {
    backgroundColor: '#e53935',
    borderRadius: 10,
    minWidth: 20,
    height: 20,
    alignItems: 'center',
    justifyContent: 'center',
    marginRight: 8,
    paddingHorizontal: 5,
  },
  badgeText: {
    color: '#fff',
    fontSize: 12,
    fontWeight: 'bold',
  },
  logoutBtn: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    marginHorizontal: 16,
    marginTop: 24,
    backgroundColor: '#fff',
    borderRadius: 12,
    paddingVertical: 14,
    borderWidth: 1,
    borderColor: '#1976d2',
  },
  logoutText: {
    color: '#1976d2',
    fontSize: 16,
    fontWeight: 'bold',
  },
  fab: {
    position: 'absolute',
    bottom: 24,
    right: 24,
    backgroundColor: '#1976d2',
    width: 56,
    height: 56,
    borderRadius: 28,
    alignItems: 'center',
    justifyContent: 'center',
    shadowColor: '#1976d2',
    shadowOpacity: 0.18,
    shadowRadius: 8,
    elevation: 4,
  },
  loading: {
    position: 'absolute',
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: 'rgba(0,0,0,0.5)',
    zIndex: 1,
  },
  center: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: '#f6f8fa',
  },
});