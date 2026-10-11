<template>
  <div class="booking-view">
    <div class="booking-container">
      <header class="booking-header">
        <p class="eyebrow">CINESTREAM CINEMA</p>
        <h1>Đặt vé xem phim</h1>
        <p>Chọn phim, rạp, suất chiếu và ghế ngồi yêu thích.</p>
      </header>

      <div v-if="pageError" class="booking-alert">{{ pageError }}</div>

      <section class="booking-panel">
        <div class="panel-heading"><span class="step-number">1</span><div><h2>Chọn phim và ngày chiếu</h2><p>Lịch chiếu mẫu được tạo theo ngày bạn chọn.</p></div></div>
        <div class="selection-grid">
          <label class="booking-field">
            <span>Phim</span>
            <select v-model="selectedMovieId" @change="refreshShowtimes">
              <option value="" disabled>Chọn phim muốn xem</option>
              <option v-for="movie in movies" :key="movie.id" :value="String(movie.id)">{{ movie.title }}</option>
            </select>
          </label>
          <label class="booking-field">
            <span>Ngày chiếu</span>
            <input v-model="selectedDate" type="date" :min="minDate" :max="maxDate" @change="refreshShowtimes" />
          </label>
        </div>
      </section>

      <section class="booking-panel" :class="{ muted: !selectedMovie }">
        <div class="panel-heading"><span class="step-number">2</span><div><h2>Chọn rạp</h2><p>Suất chiếu còn chỗ sẽ hiển thị sau khi chọn phim.</p></div></div>
        <div class="cinema-list">
          <button v-for="cinema in cinemas" :key="cinema.id" class="cinema-option" :class="{ selected: selectedCinemaId === cinema.id }" :disabled="!selectedMovie" @click="selectCinema(cinema)">
            <strong>{{ cinema.name }}</strong><span>{{ cinema.address }}</span>
          </button>
        </div>
      </section>

      <section class="booking-panel" :class="{ muted: !selectedCinemaId }">
        <div class="panel-heading"><span class="step-number">3</span><div><h2>Chọn suất chiếu</h2><p>10:00, 13:30 và 19:00 · Giá ghế thường 90.000đ, ghế VIP 130.000đ.</p></div></div>
        <div v-if="loadingShowtimes" class="booking-hint">Đang tải lịch chiếu…</div>
        <div v-else-if="selectedCinemaId && !showtimes.length" class="booking-hint">Ngày này chưa có suất chiếu khả dụng.</div>
        <div v-else class="showtime-list">
          <button v-for="showtime in showtimes" :key="showtime.id" class="showtime-option" :class="{ selected: selectedShowtime?.id === showtime.id }" :disabled="showtime.available_seats === 0" @click="selectShowtime(showtime)">
            <strong>{{ formatTime(showtime.start_at) }}</strong><span>{{ showtime.auditorium?.name }}</span><small>{{ showtime.available_seats }} ghế trống</small>
          </button>
        </div>
      </section>

      <section v-if="selectedShowtime" class="booking-panel seat-panel">
        <div class="panel-heading"><span class="step-number">4</span><div><h2>Chọn ghế ngồi</h2><p>{{ selectedShowtime.cinema?.name }} · {{ selectedShowtime.auditorium?.name }}</p></div></div>
        <div class="screen"><span>MÀN HÌNH</span></div>
        <div class="seat-legend">
          <span><i class="seat-swatch available"></i>Còn trống</span><span><i class="seat-swatch selected"></i>Đang chọn</span><span><i class="seat-swatch reserved"></i>Đã giữ/đặt</span><span><i class="seat-swatch vip"></i>VIP</span>
        </div>
        <div class="seat-map">
          <div v-for="row in seatRows" :key="row.name" class="seat-row">
            <span class="row-label">{{ row.name }}</span>
            <button v-for="seat in row.seats" :key="seat.id" type="button" class="seat" :class="[seat.status, seat.type, { chosen: selectedSeats.includes(seat.code) }]" :disabled="seat.status !== 'available'" :title="`${seat.code} · ${seat.type === 'vip' ? 'VIP' : 'Thường'}`" @click="toggleSeat(seat)">{{ seat.number }}</button>
            <span class="row-label">{{ row.name }}</span>
          </div>
        </div>
        <div class="seat-summary">
          <div><span>Ghế đã chọn</span><strong>{{ selectedSeats.length ? selectedSeats.join(', ') : 'Chưa chọn' }}</strong></div>
          <div><span>Tạm tính</span><strong>{{ formatCurrency(totalAmount) }}</strong></div>
        </div>
        <button class="btn btn-primary booking-action" :disabled="!selectedSeats.length || busy" @click="holdSeats">{{ busy ? 'Đang xử lý…' : 'Giữ ghế và tiếp tục' }}</button>
      </section>

      <section v-if="booking" class="booking-panel checkout-panel">
        <div class="panel-heading"><span class="step-number">5</span><div><h2>{{ booking.status === 'CONFIRMED' ? 'Vé đã xác nhận' : 'Thanh toán ZaloPay Sandbox' }}</h2><p>Mã vé <strong>{{ booking.code }}</strong> · Ghế {{ booking.seat_codes.join(', ') }}</p></div></div>
        <div class="checkout-summary"><span>{{ booking.showtime.movie_title }} · {{ formatTime(booking.showtime.start_at) }}</span><strong>{{ formatCurrency(booking.total_amount) }}</strong></div>
        <div v-if="booking.status === 'AWAITING_PAYMENT'" class="payment-note">Đơn giữ ghế có hiệu lực trong 10 phút. Tiếp tục để mở cổng thanh toán ZaloPay Sandbox.</div>
        <div v-if="booking.status === 'CONFIRMED'" class="ticket-success"><strong>Thanh toán thành công</strong><span>Đưa mã vé {{ booking.code }} tại quầy để nhận vé.</span></div>
        <div class="checkout-actions" v-if="booking.status === 'AWAITING_PAYMENT'">
          <button class="btn btn-primary" :disabled="busy" @click="payBooking">{{ busy ? 'Đang chuyển tới ZaloPay…' : 'Thanh toán bằng ZaloPay' }}</button>
          <button class="btn btn-secondary" :disabled="busy" @click="cancelCurrentBooking">Hủy giữ ghế</button>
        </div>
      </section>

      <section v-if="authState.isAuthenticated" class="booking-panel history-panel">
        <div class="panel-heading"><span class="step-number">✓</span><div><h2>Vé của tôi</h2><p>Các đơn đặt vé đã lưu trong tài khoản.</p></div></div>
        <p v-if="historyError" class="booking-alert">{{ historyError }}</p>
        <p v-else-if="!bookingHistory.length" class="booking-hint">Bạn chưa có đơn đặt vé nào.</p>
        <article v-for="ticket in bookingHistory" :key="ticket.code" class="history-ticket">
          <div><strong>{{ ticket.showtime.movie_title }}</strong><span>{{ ticket.showtime.cinema?.name }} · {{ ticket.showtime.auditorium?.name }}</span></div>
          <div><strong>{{ formatTime(ticket.showtime.start_at) }} · {{ formatDateTime(ticket.showtime.start_at) }}</strong><span>Ghế {{ ticket.seat_codes.join(', ') }} · Mã {{ ticket.code }}</span></div>
          <div class="history-status" :class="ticket.status.toLowerCase()"><strong>{{ statusLabel(ticket.status) }}</strong><span>{{ formatCurrency(ticket.total_amount) }}</span></div>
        </article>
      </section>

      <p v-if="actionError" class="booking-alert">{{ actionError }}</p>
      <p class="booking-disclaimer">Lịch chiếu và rạp là dữ liệu demo. Thanh toán được xử lý qua ZaloPay Sandbox.</p>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue';
import { useRoute } from 'vue-router';
import { fetchMovies } from '../services/api.js';
import { authState } from '../services/auth.js';
import { cancelBooking, createBooking, fetchCinemas, fetchMyBookings, fetchShowtimeSeats, fetchShowtimes, startZaloPayPayment, verifyBookingPayment } from '../services/booking.js';

const emit = defineEmits(['open-auth']);
const route = useRoute();
const movies = ref([]);
const cinemas = ref([]);
const showtimes = ref([]);
const seatData = ref(null);
const selectedMovieId = ref('');
const selectedDate = ref(dateString(new Date()));
const selectedCinemaId = ref(null);
const selectedShowtime = ref(null);
const selectedSeats = ref([]);
const booking = ref(null);
const bookingHistory = ref([]);
const historyError = ref('');
const loadingShowtimes = ref(false);
const busy = ref(false);
const pageError = ref('');
const actionError = ref('');

function dateString(date) {
  const local = new Date(date.getTime() - date.getTimezoneOffset() * 60000);
  return local.toISOString().slice(0, 10);
}

const minDate = dateString(new Date());
const maxDate = dateString(new Date(Date.now() + 14 * 24 * 60 * 60 * 1000));
const selectedMovie = computed(() => movies.value.find(movie => String(movie.id) === selectedMovieId.value) || null);
const seatRows = computed(() => {
  const grouped = new Map();
  for (const seat of seatData.value?.seats || []) {
    if (!grouped.has(seat.row)) grouped.set(seat.row, []);
    grouped.get(seat.row).push(seat);
  }
  return [...grouped.entries()].sort(([a], [b]) => a.localeCompare(b)).map(([name, seats]) => ({ name, seats }));
});
const totalAmount = computed(() => (seatData.value?.seats || [])
  .filter(seat => selectedSeats.value.includes(seat.code))
  .reduce((sum, seat) => sum + (seat.type === 'vip' ? selectedShowtime.value.vip_price : selectedShowtime.value.standard_price), 0));

onMounted(async () => {
  const returnedBookingCode = route.query.payment_return === '1' ? String(route.query.booking || '') : '';
  if (returnedBookingCode && authState.isAuthenticated) {
    try {
      const result = await verifyBookingPayment(returnedBookingCode);
      booking.value = result.booking;
      if (result.booking.payment_status !== 'PAID') actionError.value = 'ZaloPay chưa xác nhận thanh toán. Bạn có thể tải lại sau ít phút.';
      await refreshBookingHistory();
    } catch (error) {
      actionError.value = error.message;
    }
  }
  if (authState.isAuthenticated) refreshBookingHistory();
  try {
    const [movieList, cinemaList] = await Promise.all([fetchMovies(1), fetchCinemas()]);
    movies.value = movieList;
    cinemas.value = cinemaList;
    const queryMovie = String(route.query.movie || '');
    if (queryMovie && movieList.some(movie => String(movie.id) === queryMovie)) selectedMovieId.value = queryMovie;
    if (cinemaList.length) selectedCinemaId.value = cinemaList[0].id;
    if (selectedMovie.value && selectedCinemaId.value) await refreshShowtimes();
  } catch (error) {
    pageError.value = error.message;
  }
});

watch(() => authState.token, (token) => {
  if (token) refreshBookingHistory();
  else bookingHistory.value = [];
});

async function refreshBookingHistory() {
  historyError.value = '';
  try {
    const history = await fetchMyBookings();
    bookingHistory.value = await Promise.all(history.map(async (ticket) => {
      if (ticket.status !== 'AWAITING_PAYMENT') return ticket;
      try {
        const result = await verifyBookingPayment(ticket.code);
        const updated = result.booking || ticket;
        if (booking.value?.code === ticket.code) booking.value = updated;
        return updated;
      } catch {
        return ticket;
      }
    }));
  } catch (error) {
    historyError.value = error.message;
  }
}

async function refreshShowtimes() {
  selectedShowtime.value = null;
  seatData.value = null;
  selectedSeats.value = [];
  booking.value = null;
  actionError.value = '';
  if (!selectedMovie.value || !selectedCinemaId.value) {
    showtimes.value = [];
    return;
  }
  loadingShowtimes.value = true;
  try {
    showtimes.value = await fetchShowtimes(selectedMovie.value, selectedDate.value, selectedCinemaId.value);
  } catch (error) {
    actionError.value = error.message;
  } finally {
    loadingShowtimes.value = false;
  }
}

async function selectCinema(cinema) {
  selectedCinemaId.value = cinema.id;
  await refreshShowtimes();
}

async function selectShowtime(showtime) {
  selectedShowtime.value = showtime;
  selectedSeats.value = [];
  booking.value = null;
  actionError.value = '';
  await refreshSeatMap();
}

async function refreshSeatMap() {
  try {
    seatData.value = await fetchShowtimeSeats(selectedShowtime.value.id);
  } catch (error) {
    actionError.value = error.message;
  }
}

function toggleSeat(seat) {
  const index = selectedSeats.value.indexOf(seat.code);
  if (index >= 0) selectedSeats.value.splice(index, 1);
  else if (selectedSeats.value.length >= 8) actionError.value = 'Mỗi đơn đặt tối đa 8 ghế.';
  else {
    actionError.value = '';
    selectedSeats.value.push(seat.code);
  }
}

async function holdSeats() {
  if (!authState.isAuthenticated) {
    emit('open-auth');
    actionError.value = 'Đăng nhập để giữ ghế và nhận vé.';
    return;
  }
  busy.value = true;
  actionError.value = '';
  try {
    const result = await createBooking(selectedShowtime.value.id, selectedSeats.value);
    booking.value = result.booking;
    selectedSeats.value = [];
    await refreshSeatMap();
    await refreshBookingHistory();
  } catch (error) {
    actionError.value = error.message;
    await refreshSeatMap();
  } finally {
    busy.value = false;
  }
}

async function payBooking() {
  busy.value = true;
  actionError.value = '';
  try {
    const result = await startZaloPayPayment(booking.value.code);
    if (result.booking) {
      booking.value = result.booking;
      await refreshSeatMap();
      await refreshBookingHistory();
      return;
    }
    if (!result.order_url) throw new Error('ZaloPay không trả về đường dẫn thanh toán.');
    window.location.assign(result.order_url);
  } catch (error) {
    actionError.value = error.message;
  } finally {
    busy.value = false;
  }
}

async function cancelCurrentBooking() {
  busy.value = true;
  actionError.value = '';
  try {
    await cancelBooking(booking.value.code);
    booking.value = null;
    selectedSeats.value = [];
    await refreshSeatMap();
    await refreshBookingHistory();
  } catch (error) {
    actionError.value = error.message;
  } finally {
    busy.value = false;
  }
}

function formatTime(value) {
  return new Intl.DateTimeFormat('vi-VN', { hour: '2-digit', minute: '2-digit', timeZone: 'Asia/Ho_Chi_Minh' }).format(new Date(value));
}

function formatDateTime(value) {
  return new Intl.DateTimeFormat('vi-VN', { day: '2-digit', month: '2-digit', year: 'numeric', timeZone: 'Asia/Ho_Chi_Minh' }).format(new Date(value));
}

function statusLabel(status) {
  return ({ CONFIRMED: 'Đã thanh toán', AWAITING_PAYMENT: 'Chờ thanh toán', EXPIRED: 'Đã hết hạn', CANCELLED: 'Đã hủy' })[status] || status;
}

function formatCurrency(amount) {
  return new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND', maximumFractionDigits: 0 }).format(amount || 0);
}
</script>

<style scoped>
.booking-view { min-height: 80vh; padding: 116px 24px 64px; color: var(--text-primary); }
.booking-container { width: min(100%, 980px); margin: 0 auto; }
.booking-header { margin-bottom: 28px; }
.booking-header h1 { margin: 4px 0 8px; font-family: var(--font-heading); font-size: clamp(2rem, 5vw, 3rem); }
.booking-header p:last-child, .panel-heading p { color: var(--text-secondary); margin: 0; }
.eyebrow { color: var(--accent-cyan); font-size: .75rem; font-weight: 800; letter-spacing: .2em; }
.booking-panel { margin: 18px 0; padding: 24px; border: 1px solid var(--border-subtle); border-radius: 18px; background: var(--bg-card); }
.booking-panel.muted { opacity: .65; }
.panel-heading { display: flex; align-items: flex-start; gap: 14px; margin-bottom: 20px; }
.panel-heading h2 { margin: 0 0 5px; font-size: 1.1rem; }
.panel-heading p { font-size: .88rem; }
.step-number { display: grid; place-items: center; flex: 0 0 34px; height: 34px; border-radius: 50%; color: #fff; font-weight: 800; background: var(--accent-gradient); }
.selection-grid { display: grid; grid-template-columns: 2fr 1fr; gap: 14px; }
.booking-field { display: grid; gap: 8px; color: var(--text-secondary); font-size: .85rem; font-weight: 700; }
.booking-field select, .booking-field input { min-width: 0; padding: 12px 14px; border: 1px solid var(--border-subtle); border-radius: 10px; color: var(--text-primary); background: #111827; }
.booking-field option { background: #111827; }
.cinema-list { display: grid; grid-template-columns: repeat(auto-fit, minmax(220px, 1fr)); gap: 12px; }
.cinema-option, .showtime-option { display: grid; gap: 5px; padding: 14px; text-align: left; color: var(--text-primary); background: rgba(255,255,255,.035); border: 1px solid var(--border-subtle); border-radius: 12px; cursor: pointer; transition: .2s ease; }
.cinema-option span, .showtime-option span, .showtime-option small { color: var(--text-secondary); font-size: .8rem; }
.cinema-option.selected, .showtime-option.selected { border-color: var(--accent-cyan); background: rgba(0,242,254,.09); }
.cinema-option:disabled { cursor: not-allowed; }
.showtime-list { display: flex; flex-wrap: wrap; gap: 10px; }
.showtime-option { min-width: 150px; text-align: center; }
.showtime-option strong { font-size: 1.2rem; }
.showtime-option:disabled { opacity: .45; cursor: not-allowed; }
.booking-hint { color: var(--text-secondary); font-size: .9rem; }
.screen { width: min(80%, 560px); margin: 12px auto 28px; padding: 10px; text-align: center; font-size: .68rem; font-weight: 800; letter-spacing: .3em; color: var(--accent-cyan); border-top: 3px solid var(--accent-cyan); border-radius: 50% 50% 0 0 / 22px 22px 0 0; box-shadow: 0 -8px 22px rgba(0,242,254,.12); }
.seat-legend { display: flex; flex-wrap: wrap; justify-content: center; gap: 14px; margin-bottom: 18px; color: var(--text-secondary); font-size: .75rem; }
.seat-legend span { display: inline-flex; align-items: center; gap: 6px; }
.seat-swatch { width: 14px; height: 14px; display: inline-block; border-radius: 4px; background: #253148; border: 1px solid #42506a; }
.seat-swatch.selected, .seat.chosen { background: var(--accent-cyan); border-color: var(--accent-cyan); color: #041014; }
.seat-swatch.reserved, .seat.reserved, .seat.booked { background: #4b5563; color: #9ca3af; border-color: #4b5563; }
.seat-swatch.vip, .seat.vip:not(.chosen):not(.reserved):not(.booked) { border-color: #b88a27; color: #f8ce69; }
.seat-map { width: fit-content; max-width: 100%; display: grid; gap: 8px; margin: 0 auto; overflow-x: auto; }
.seat-row { display: flex; align-items: center; gap: 6px; }
.row-label { width: 18px; text-align: center; color: var(--text-secondary); font-size: .75rem; }
.seat { width: 32px; height: 30px; padding: 0; color: #d7e2f3; border: 1px solid #42506a; border-radius: 7px 7px 4px 4px; background: #253148; cursor: pointer; font-size: .7rem; transition: .15s ease; }
.seat:hover:not(:disabled) { transform: translateY(-2px); border-color: var(--accent-cyan); }
.seat:disabled { cursor: not-allowed; }
.seat-summary, .checkout-summary { display: flex; justify-content: space-between; gap: 14px; margin: 24px 0 14px; padding: 16px; border-radius: 12px; background: rgba(255,255,255,.045); }
.seat-summary div { display: grid; gap: 5px; }
.seat-summary span, .checkout-summary span { color: var(--text-secondary); font-size: .82rem; }
.seat-summary strong { overflow-wrap: anywhere; }
.booking-action { display: block; margin-left: auto; }
.checkout-panel { border-color: rgba(0,242,254,.35); }
.payment-note, .booking-disclaimer { color: var(--text-secondary); font-size: .86rem; line-height: 1.6; }
.checkout-actions { display: flex; flex-wrap: wrap; gap: 10px; margin-top: 16px; }
.ticket-success { display: grid; gap: 5px; padding: 14px; color: #a7f3d0; background: rgba(16,185,129,.1); border-radius: 10px; }
.history-ticket { display: grid; grid-template-columns: 1.2fr 1fr auto; align-items: center; gap: 16px; padding: 14px 0; border-top: 1px solid var(--border-subtle); }
.history-ticket > div { display: grid; gap: 5px; }
.history-ticket span { color: var(--text-secondary); font-size: .8rem; }
.history-status.confirmed { color: #6ee7b7; }
.history-status.awaiting_payment { color: #fcd34d; }
.history-status.expired, .history-status.cancelled { color: #fca5a5; }
.booking-alert { padding: 13px 16px; border: 1px solid rgba(229,9,20,.5); border-radius: 10px; color: #fecaca; background: rgba(229,9,20,.12); }
.booking-disclaimer { text-align: center; margin-top: 24px; }
@media (max-width: 640px) { .booking-view { padding: 100px 14px 44px; } .booking-panel { padding: 18px 14px; } .selection-grid { grid-template-columns: 1fr; } .seat { width: 27px; height: 27px; } .seat-row { gap: 4px; } .history-ticket { grid-template-columns: 1fr; gap: 8px; } }
</style>
