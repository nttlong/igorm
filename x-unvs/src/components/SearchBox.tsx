import { useState } from 'react';
import { FiMenu, FiBell, FiSun, FiMoon, FiLogOut,FiSearch } from "react-icons/fi";
function SearchBox() {
  const [searchTerm, setSearchTerm] = useState('');

  const handleSearch = () => {
    if (searchTerm.trim()) {
      console.log('Searching for:', searchTerm);
      // Thêm logic tìm kiếm tại đây
    }
  };

  return (
    <div className="flex">
      <input
        type="text"
        value={searchTerm}
        onChange={(e) => setSearchTerm(e.target.value)}
        placeholder="Enter search term..."
        className="flex-1 p-2 border border-gray-300 rounded-l-md focus:outline-none focus:ring-2 focus:ring-blue-500"
      />
      <button
        onClick={handleSearch}
        className="bg-blue-500 text-white p-2 rounded-r-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500"
      >
        <FiSearch size={16} />
      </button>
    </div>
  );
}

export default SearchBox;