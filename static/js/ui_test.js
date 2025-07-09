/**
 * UI Tests for Pack Calculator
 * 
 * This file contains tests for the Pack Calculator UI.
 * These tests can be run in a browser using a testing framework like Jest.
 */

// Mock fetch function for testing
const mockFetch = (response) => {
  return jest.fn().mockImplementation(() => 
    Promise.resolve({
      ok: true,
      json: () => Promise.resolve(response)
    })
  );
};

describe('Pack Calculator UI', () => {
  // Save the original fetch
  const originalFetch = global.fetch;
  
  beforeEach(() => {
    // Create a mock DOM for testing
    document.body.innerHTML = `
      <div id="pack-sizes-container"></div>
      <input id="new-pack-size" type="number" value="250">
      <button id="add-pack-size">Add</button>
      <div id="pack-sizes-error" class="error hidden"></div>
      
      <input id="items-ordered" type="number" value="501">
      <button id="calculate-btn">Calculate</button>
      <div id="calculation-result" class="result hidden">
        <span id="result-items-ordered">0</span>
        <span id="result-total-items">0</span>
        <span id="result-total-packs">0</span>
        <div id="result-packs"></div>
      </div>
      <div id="calculation-error" class="error hidden"></div>
    `;
    
    // Mock fetch
    global.fetch = jest.fn();
  });
  
  afterEach(() => {
    // Restore the original fetch
    global.fetch = originalFetch;
    
    // Clear all mocks
    jest.clearAllMocks();
  });
  
  test('should load pack sizes on page load', async () => {
    // Mock the response
    const mockResponse = {
      data: {
        pack_sizes: {
          items: [
            { id: '1', size: 250 },
            { id: '2', size: 500 },
            { id: '3', size: 1000 },
            { id: '4', size: 2000 },
            { id: '5', size: 5000 }
          ]
        }
      }
    };
    
    global.fetch = mockFetch(mockResponse);
    
    // Call the loadPackSizes function
    await loadPackSizes();
    
    // Check that fetch was called with the correct arguments
    expect(global.fetch).toHaveBeenCalledWith('/graphql', expect.any(Object));
    
    // Check that the pack sizes were rendered
    const packSizesContainer = document.getElementById('pack-sizes-container');
    expect(packSizesContainer.children.length).toBe(5);
    
    // Check that the first pack size is rendered correctly
    const firstPackSize = packSizesContainer.children[0];
    expect(firstPackSize.textContent).toContain('250 items');
  });
  
  test('should add a new pack size', async () => {
    // Mock the response
    const mockResponse = {
      data: {
        create_pack_size: {
          id: '6',
          size: 300
        }
      }
    };
    
    global.fetch = mockFetch(mockResponse);
    
    // Set the input value
    document.getElementById('new-pack-size').value = '300';
    
    // Call the addPackSize function
    await addPackSize();
    
    // Check that fetch was called with the correct arguments
    expect(global.fetch).toHaveBeenCalledWith('/graphql', expect.any(Object));
    
    // Check that the input was cleared
    expect(document.getElementById('new-pack-size').value).toBe('');
    
    // Check that loadPackSizes was called
    expect(global.fetch).toHaveBeenCalledTimes(2);
  });
  
  test('should calculate packs for an order', async () => {
    // Mock the response
    const mockResponse = {
      data: {
        calculate_packs: {
          items_ordered: 501,
          total_items: 750,
          total_packs: 2,
          packs: [
            { pack_size: 500, quantity: 1 },
            { pack_size: 250, quantity: 1 }
          ]
        }
      }
    };
    
    global.fetch = mockFetch(mockResponse);
    
    // Set the input value
    document.getElementById('items-ordered').value = '501';
    
    // Call the calculatePacks function
    await calculatePacks();
    
    // Check that fetch was called with the correct arguments
    expect(global.fetch).toHaveBeenCalledWith('/graphql', expect.any(Object));
    
    // Check that the result is displayed
    const calculationResult = document.getElementById('calculation-result');
    expect(calculationResult.classList.contains('hidden')).toBe(false);
    
    // Check that the result values are correct
    expect(document.getElementById('result-items-ordered').textContent).toBe('501');
    expect(document.getElementById('result-total-items').textContent).toBe('750');
    expect(document.getElementById('result-total-packs').textContent).toBe('2');
    
    // Check that the packs are rendered
    const resultPacks = document.getElementById('result-packs');
    expect(resultPacks.children.length).toBe(2);
    
    // Check that the first pack is rendered correctly
    const firstPack = resultPacks.children[0];
    expect(firstPack.textContent).toContain('1 × 500 item pack');
    
    // Check that the second pack is rendered correctly
    const secondPack = resultPacks.children[1];
    expect(secondPack.textContent).toContain('1 × 250 item pack');
  });
  
  test('should handle errors when calculating packs', async () => {
    // Mock the response
    const mockResponse = {
      errors: [
        { message: 'items ordered must be greater than 0' }
      ]
    };
    
    global.fetch = mockFetch(mockResponse);
    
    // Set the input value
    document.getElementById('items-ordered').value = '0';
    
    // Call the calculatePacks function
    await calculatePacks();
    
    // Check that fetch was called with the correct arguments
    expect(global.fetch).toHaveBeenCalledWith('/graphql', expect.any(Object));
    
    // Check that the error is displayed
    const calculationError = document.getElementById('calculation-error');
    expect(calculationError.classList.contains('hidden')).toBe(false);
    expect(calculationError.textContent).toContain('items ordered must be greater than 0');
    
    // Check that the result is hidden
    const calculationResult = document.getElementById('calculation-result');
    expect(calculationResult.classList.contains('hidden')).toBe(true);
  });
});
