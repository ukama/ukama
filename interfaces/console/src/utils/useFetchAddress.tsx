import { useState, useCallback } from 'react';

interface UseFetchAddressResult {
  address: string;
  isLoading: boolean;
  error: string | null;
  fetchAddress: (lat: number, lng: number) => Promise<void>;
}

export const useFetchAddress = (): UseFetchAddressResult => {
  const [address, setAddress] = useState<string>('');
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  const fetchAddress = useCallback(async (lat: number, lng: number) => {
    setIsLoading(true);
    setError(null);

    try {
      const response = await fetch(
        `https://nominatim.openstreetmap.org/reverse?format=json&lat=${lat || 37.7749}&lon=${lng || -122.4194}`,
        {
          cache: 'force-cache',
        },
      );

      if (!response.ok) {
        throw new Error('Failed to fetch address');
      }

      const data = await response.json();
      setAddress(data.display_name || 'Location not found');
    } catch (err) {
      setError('Error fetching address');
      setAddress('Location not found');
    } finally {
      setIsLoading(false);
    }
  }, []);

  return { address, isLoading, error, fetchAddress };
};
