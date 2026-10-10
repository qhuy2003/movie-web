const API_BASE_URL = 'http://localhost:8080/movies';

// High Quality Fallback & Curated Movie Database
export const MOCK_MOVIES = [
  {
    id: 101,
    title: "Dune: Hành Tinh Cát - Phần 2",
    origin_name: "Dune: Part Two",
    year: 2024,
    rating: 8.8,
    quality: "4K UHD",
    duration: "166 phút",
    director: "Denis Villeneuve",
    cast: ["Timothée Chalamet", "Zendaya", "Rebecca Ferguson", "Javier Bardem"],
    genres: ["Hành Động", "Viễn Tưởng", "Phiêu Lưu"],
    poster_url: "https://images.unsplash.com/photo-1534447677768-be436bb09401?w=600&auto=format&fit=crop&q=80",
    thumb_url: "https://images.unsplash.com/photo-1518709268805-4e9042af9f23?w=1200&auto=format&fit=crop&q=80",
    description: "Paul Atreides tái hợp với Chani và người Fremen trong chuyến hành trình trả thù những kẻ đã hủy hoại gia đình anh. Anh phải lựa chọn giữa tình yêu của đời mình và số phận của cả vũ trụ.",
    video_url: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/TearsOfSteel.mp4",
    episodes: [
      { id: 1, name: "Bản Chiếu Rạp (4K Vietsub)", url: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/TearsOfSteel.mp4" },
      { id: 2, name: "Bản Mở Rộng IMAX (Thuyết Minh)", url: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/Sintel.mp4" }
    ],
    is_featured: true,
    is_trending: true
  },
  {
    id: 102,
    title: "Oppenheimer",
    origin_name: "Oppenheimer",
    year: 2023,
    rating: 8.9,
    quality: "4K UHD",
    duration: "180 phút",
    director: "Christopher Nolan",
    cast: ["Cillian Murphy", "Emily Blunt", "Matt Damon", "Robert Downey Jr."],
    genres: ["Chính Kịch", "Lịch Sử", "Tiểu Sử"],
    poster_url: "https://images.unsplash.com/photo-1440404653325-ab127d49abc1?w=600&auto=format&fit=crop&q=80",
    thumb_url: "https://images.unsplash.com/photo-1509198397868-475647b2a1e5?w=1200&auto=format&fit=crop&q=80",
    description: "Câu chuyện về nhà vật lý lý thuyết J. Robert Oppenheimer, người lãnh đạo Dự án Manhattan chế tạo ra bom nguyên tử đầu tiên cho loài người.",
    video_url: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/BigBuckBunny.mp4",
    episodes: [
      { id: 1, name: "Full HD Vietsub", url: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/BigBuckBunny.mp4" }
    ],
    is_featured: true,
    is_trending: true
  },
  {
    id: 103,
    title: "Spider-Man: Du Hành Vũ Trụ Nhện",
    origin_name: "Spider-Man: Across the Spider-Verse",
    year: 2023,
    rating: 8.7,
    quality: "Full HD",
    duration: "140 phút",
    director: "Joaquim Dos Santos",
    cast: ["Shameik Moore", "Hailee Steinfeld", "Oscar Isaac"],
    genres: ["Hoạt Hình", "Hành Động", "Viễn Tưởng"],
    poster_url: "https://images.unsplash.com/photo-1635805737707-575885ab0820?w=600&auto=format&fit=crop&q=80",
    thumb_url: "https://images.unsplash.com/photo-1607604276583-eef5d076aa5f?w=1200&auto=format&fit=crop&q=80",
    description: "Miles Morales phóng qua Đa vũ trụ, nơi anh gặp một nhóm Người Nhện chịu trách nhiệm bảo vệ sự tồn tại của nó. Nhưng khi các anh hùng xung đột về cách xử lý mối đe dọa mới, Miles phải định nghĩa lại ý nghĩa của việc làm một anh hùng.",
    video_url: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ElephantsDream.mp4",
    episodes: [
      { id: 1, name: "Full HD Vietsub", url: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ElephantsDream.mp4" }
    ],
    is_featured: false,
    is_trending: true
  },
  {
    id: 104,
    title: "Inception: Kẻ Đánh Cắp Giấc Mơ",
    origin_name: "Inception",
    year: 2010,
    rating: 8.8,
    quality: "4K UHD",
    duration: "148 phút",
    director: "Christopher Nolan",
    cast: ["Leonardo DiCaprio", "Joseph Gordon-Levitt", "Elliot Page"],
    genres: ["Hành Động", "Viễn Tưởng", "Gây Cấn"],
    poster_url: "https://images.unsplash.com/photo-1579783902614-a3fb3927b675?w=600&auto=format&fit=crop&q=80",
    thumb_url: "https://images.unsplash.com/photo-1518709268805-4e9042af9f23?w=1200&auto=format&fit=crop&q=80",
    description: "Dom Cobb là một kẻ cắp chuyên nghiệp, chuyên đánh cắp bí mật của người khác bằng cách xâm nhập vào giấc mơ của họ.",
    video_url: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerBlazes.mp4",
    episodes: [
      { id: 1, name: "Full HD Vietsub", url: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerBlazes.mp4" }
    ],
    is_featured: true,
    is_trending: false
  },
  {
    id: 105,
    title: "Interstellar: Hố Đen Vũ Trụ",
    origin_name: "Interstellar",
    year: 2014,
    rating: 8.7,
    quality: "4K UHD",
    duration: "169 phút",
    director: "Christopher Nolan",
    cast: ["Matthew McConaughey", "Anne Hathaway", "Jessica Chastain"],
    genres: ["Viễn Tưởng", "Phiêu Lưu", "Chính Kịch"],
    poster_url: "https://images.unsplash.com/photo-1451187580459-43490279c0fa?w=600&auto=format&fit=crop&q=80",
    thumb_url: "https://images.unsplash.com/photo-1446776811953-b23d57bd21aa?w=1200&auto=format&fit=crop&q=80",
    description: "Nhóm các nhà hành tinh học du hành qua một lỗ sâu vũ trụ gần Sao Thổ nhằm tìm kiếm một hành tinh có thể sinh sống mới cho nhân loại.",
    video_url: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerEscapes.mp4",
    episodes: [
      { id: 1, name: "Full HD Vietsub", url: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerEscapes.mp4" }
    ],
    is_featured: false,
    is_trending: true
  },
  {
    id: 106,
    title: "Cyberpunk: Edgerunners",
    origin_name: "Cyberpunk: Edgerunners",
    year: 2022,
    rating: 8.6,
    quality: "Full HD",
    duration: "10 Tập",
    director: "Hiroyuki Imaishi",
    cast: ["KENN", "Aoi Yuuki", "Hiroki Touchi"],
    genres: ["Hoạt Hình", "Hành Động", "Cyberpunk"],
    poster_url: "https://images.unsplash.com/photo-1578632767115-351597cf2477?w=600&auto=format&fit=crop&q=80",
    thumb_url: "https://images.unsplash.com/photo-1542751371-adc38448a05e?w=1200&auto=format&fit=crop&q=80",
    description: "Một đứa trẻ đường phố cố gắng sống sót trong một thành phố tương lai bị ám ảnh bởi công nghệ và chỉnh sửa cơ thể.",
    video_url: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerFun.mp4",
    episodes: [
      { id: 1, name: "Tập 1: CODE", url: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerFun.mp4" },
      { id: 2, name: "Tập 2: LIKE A BOY", url: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerJoyrides.mp4" },
      { id: 3, name: "Tập 3: SMOOTH CRIMINAL", url: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerMeltdowns.mp4" }
    ],
    is_featured: false,
    is_trending: true
  },
  {
    id: 107,
    title: "Avatar: Dòng Chảy Của Nước",
    origin_name: "Avatar: The Way of Water",
    year: 2022,
    rating: 7.7,
    quality: "4K UHD",
    duration: "192 phút",
    director: "James Cameron",
    cast: ["Sam Worthington", "Zoe Saldana", "Sigourney Weaver"],
    genres: ["Hành Động", "Viễn Tưởng", "Phiêu Lưu"],
    poster_url: "https://images.unsplash.com/photo-1518709268805-4e9042af9f23?w=600&auto=format&fit=crop&q=80",
    thumb_url: "https://images.unsplash.com/photo-1507525428034-b723cf961d3e?w=1200&auto=format&fit=crop&q=80",
    description: "Jake Sully sống cùng gia đình mới của anh trên hành tinh Pandora. Khi một mối đe dọa quen thuộc trở lại để hoàn thành công việc trước đó, Jake phải làm việc với Neytiri và quân đội của bộ tộc Na'vi để bảo vệ hành tinh của họ.",
    video_url: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/Sintel.mp4",
    episodes: [
      { id: 1, name: "Full HD Vietsub", url: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/Sintel.mp4" }
    ],
    is_featured: false,
    is_trending: false
  },
  {
    id: 108,
    title: "Chú Thuật Hồi Chiến: Chú Thuật Sư 0",
    origin_name: "Jujutsu Kaisen 0",
    year: 2021,
    rating: 8.5,
    quality: "Full HD",
    duration: "105 phút",
    director: "Sunghoo Park",
    cast: ["Megumi Ogata", "Kana Hanazawa", "Yuuichi Nakamura"],
    genres: ["Hoạt Hình", "Hành Động", "Kỳ Ảo"],
    poster_url: "https://images.unsplash.com/photo-1563089145-599997674d42?w=600&auto=format&fit=crop&q=80",
    thumb_url: "https://images.unsplash.com/photo-1618005182384-a83a8bd57fbe?w=1200&auto=format&fit=crop&q=80",
    description: "Yuta Okkotsu, một học sinh trung học bị ám bởi linh hồn người bạn thơ ấu Rika, đăng ký học tại Chú thuật Trường trung học Tokyo dưới sự hướng dẫn của Satoru Gojo.",
    video_url: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/SubaruOutback2015.mp4",
    episodes: [
      { id: 1, name: "Full HD Vietsub", url: "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/SubaruOutback2015.mp4" }
    ],
    is_featured: false,
    is_trending: true
  }
];

export async function fetchMovies(page = 1) {
  try {
    const res = await fetch(`${API_BASE_URL}?page=${page}`);
    if (!res.ok) throw new Error(`Go Backend error status ${res.status}`);
    const data = await res.json();
    
    if (Array.isArray(data) && data.length > 0) {
      // Enrich Go backend items with fallback visual assets if missing
      return data.map((item, index) => {
        const fallback = MOCK_MOVIES[index % MOCK_MOVIES.length];
        return {
          id: item.id || fallback.id,
          title: item.title || fallback.title,
          origin_name: item.origin_name || fallback.origin_name,
          slug: item.slug,
          year: item.year || fallback.year,
          rating: item.rating || fallback.rating,
          genres: (item.genres && item.genres.length) ? item.genres : fallback.genres,
          poster_url: item.poster_url || fallback.poster_url,
          thumb_url: item.thumb_url || fallback.thumb_url,
          quality: fallback.quality,
          duration: fallback.duration,
          director: fallback.director,
          cast: fallback.cast,
          description: fallback.description,
          video_url: fallback.video_url,
          episodes: fallback.episodes,
          is_featured: index === 0,
          is_trending: item.rating >= 8.0
        };
      });
    }
  } catch (err) {
    console.warn('⚠️ Go backend request failed or offline. Using rich fallback movie database:', err.message);
  }

  // Graceful Fallback
  return MOCK_MOVIES;
}

export async function fetchMovieById(id) {
  try {
    const res = await fetch(`${API_BASE_URL}/${id}`);
    if (res.ok) {
      const item = await res.json();
      const fallback = MOCK_MOVIES.find(m => m.id === item.id) || MOCK_MOVIES[0];
      return {
        ...fallback,
        id: item.id,
        title: item.title || fallback.title,
        year: item.year || fallback.year,
        rating: item.rating || fallback.rating
      };
    }
  } catch (e) {
    console.warn('Could not fetch movie by ID from backend, fallbacking to local store.');
  }

  return MOCK_MOVIES.find(m => m.id === Number(id)) || MOCK_MOVIES[0];
}

export const GENRES_LIST = [
  "Tất Cả",
  "Hành Động",
  "Hành Động & Phiêu Lưu",
  "Viễn Tưởng",
  "Khoa Học Viễn Tưởng",
  "Chính Kịch",
  "Hoạt Hình",
  "Hài",
  "Giả Tưởng",
  "Giật Gân",
  "Bí Ẩn",
  "Chiến Tranh",
  "Cổ Trang",
  "Drama",
  "Gia Đình",
  "Học Đường",
  "Kiếm Hiệp",
  "Kinh Dị",
  "Lãng Mạn",
  "Phiêu Lưu",
  "Tiên Hiệp",
  "Tình Yêu Ngọt Ngào",
  "Tội Phạm",
  "Võ Thuật",
  "Xã Hội Đen",
  "Sci-Fi & Fantasy"
];

export async function fetchGenres() {
  try {
    const res = await fetch('http://localhost:8080/genres');
    if (res.ok) {
      const data = await res.json();
      if (Array.isArray(data) && data.length > 0) {
        return ["Tất Cả", ...data.map(g => g.name)];
      }
    }
  } catch (e) {
    console.warn('Go backend /genres unavailable, attempting direct vsmov API...');
  }

  try {
    const res = await fetch('https://vsmov.com/api/the-loai');
    if (res.ok) {
      const payload = await res.json();
      if (payload && payload.data && Array.isArray(payload.data.items)) {
        return ["Tất Cả", ...payload.data.items.map(g => g.name)];
      }
    }
  } catch (e) {
    console.warn('vsmov genre API direct fetch failed, using fallback GENRES_LIST');
  }

  return GENRES_LIST;
}

export async function fetchMovieDetailBySlug(slug) {
  try {
    const res = await fetch(`http://localhost:8080/movies/detail/${slug}`, { cache: 'no-store' });
    if (res.ok) {
      const data = await res.json();
      if (data && data.movie) {
        return data;
      }
    }
  } catch (e) {
    console.warn('Go backend /movies/detail unavailable, trying direct vsmov API...');
  }

  try {
    const res = await fetch(`https://vsmov.com/api/phim/${slug}`, { cache: 'no-store' });
    if (res.ok) {
      return await res.json();
    }
  } catch (e) {
    console.warn('vsmov movie detail API direct fetch failed:', e);
  }

  return null;
}

export async function fetchActors(query = '') {
  try {
    const url = query ? `http://localhost:8080/actors?q=${encodeURIComponent(query)}` : 'http://localhost:8080/actors';
    const res = await fetch(url);
    if (res.ok) {
      const data = await res.json();
      if (Array.isArray(data)) return data;
    }
  } catch (e) {
    console.warn('Go backend /actors unavailable, trying direct vsmov API...');
  }

  try {
    const res = await fetch('https://vsmov.com/api/dien-vien');
    if (res.ok) {
      const payload = await res.json();
      if (payload && payload.data && Array.isArray(payload.data.items)) {
        let items = payload.data.items;
        if (query) {
          const q = query.toLowerCase();
          items = items.filter(a => a.name.toLowerCase().includes(q) || a.slug.toLowerCase().includes(q));
        }
        return items.slice(0, 100);
      }
    }
  } catch (e) {
    console.warn('vsmov actors API direct fetch failed:', e);
  }

  return [];
}

export async function fetchCountries() {
  try {
    const res = await fetch('http://localhost:8080/countries');
    if (res.ok) {
      const data = await res.json();
      if (Array.isArray(data) && data.length > 0) return data;
    }
  } catch (e) {
    console.warn('Go backend /countries unavailable, trying direct vsmov API...');
  }

  try {
    const res = await fetch('https://vsmov.com/api/quoc-gia');
    if (res.ok) {
      const payload = await res.json();
      if (payload && payload.data && Array.isArray(payload.data.items)) {
        return payload.data.items;
      }
    }
  } catch (e) {
    console.warn('vsmov countries API direct fetch failed:', e);
  }

  return [
    { _id: 17, name: "Nhật Bản", slug: "nhat-ban" },
    { _id: 10, name: "Hàn Quốc", slug: "han-quoc" },
    { _id: 13, name: "Trung Quốc", slug: "trung-quoc" },
    { _id: 7, name: "Âu Mỹ", slug: "au-my" },
    { _id: 9, name: "Việt Nam", slug: "viet-nam" },
    { _id: 14, name: "Thái Lan", slug: "thai-lan" }
  ];
}

export async function fetchMoviesByCountry(slug, filters = {}) {
  const params = new URLSearchParams();
  if (filters.page) params.set('page', filters.page);
  if (filters.limit) params.set('limit', filters.limit);
  if (filters.year) params.set('year', filters.year);
  if (filters.type) params.set('type', filters.type);
  if (filters.status) params.set('status', filters.status);

  const queryString = params.toString() ? `?${params.toString()}` : '';

  try {
    const res = await fetch(`http://localhost:8080/countries/${slug}${queryString}`);
    if (res.ok) {
      return await res.json();
    }
  } catch (e) {
    console.warn('Go backend /countries/:slug unavailable, trying direct vsmov API...');
  }

  try {
    const res = await fetch(`https://vsmov.com/api/quoc-gia/${slug}${queryString}`);
    if (res.ok) {
      return await res.json();
    }
  } catch (e) {
    console.warn('vsmov movies by country direct fetch failed:', e);
  }

  return null;
}

export async function fetchMoviesByCatalog(slug = 'subteam', filters = {}) {
  const params = new URLSearchParams();
  if (filters.page) params.set('page', filters.page);
  if (filters.limit) params.set('limit', filters.limit);
  if (filters.year) params.set('year', filters.year);
  if (filters.type) params.set('type', filters.type);
  if (filters.status) params.set('status', filters.status);

  const queryString = params.toString() ? `?${params.toString()}` : '';

  try {
    const res = await fetch(`http://localhost:8080/catalog/${slug}${queryString}`);
    if (res.ok) {
      return await res.json();
    }
  } catch (e) {
    console.warn('Go backend /catalog/:slug unavailable, trying direct vsmov API...');
  }

  try {
    const res = await fetch(`https://vsmov.com/api/danh-sach/${slug}${queryString}`);
    if (res.ok) {
      return await res.json();
    }
  } catch (e) {
    console.warn('vsmov movies by catalog direct fetch failed:', e);
  }

  return null;
}

export async function fetchYears() {
  try {
    const res = await fetch('http://localhost:8080/years');
    if (!res.ok) throw new Error(`HTTP ${res.status}`);
    const payload = await res.json();
    const items = payload?.data?.items || payload?.items || payload;
    if (Array.isArray(items) && items.length) return items;
  } catch (e) {
    console.warn('Go backend /years unavailable, using generated year list:', e.message);
  }
  const currentYear = new Date().getFullYear();
  return Array.from({ length: currentYear - 1979 }, (_, index) => ({
    name: String(currentYear - index), slug: String(currentYear - index)
  }));
}

export async function fetchMoviesByYear(year, filters = {}) {
  const params = new URLSearchParams();
  if (filters.page) params.set('page', filters.page);
  if (filters.limit) params.set('limit', filters.limit);
  if (filters.type) params.set('type', filters.type);
  if (filters.status) params.set('status', filters.status);
  const query = params.toString() ? `?${params}` : '';
  const urls = [
    `http://localhost:8080/years/${encodeURIComponent(year)}${query}`,
    `https://vsmov.com/api/nam/${encodeURIComponent(year)}${query}`
  ];
  for (const url of urls) {
    try {
      const res = await fetch(url);
      if (res.ok) return await res.json();
    } catch (e) {
      console.warn(`Year API unavailable: ${url}`, e.message);
    }
  }
  return null;
}
