import pytest
import json
from app.solar_tool.schemas.solar_tool import Site
from unittest.mock import MagicMock, Mock, mock_open, patch
from app.solar_tool.services.solar_tool import SolarTool
from math import radians

class TestSolarTool:

    @pytest.fixture(autouse=True)
    def setup(self):
        self.instance = SolarTool()
        
    @patch('app.solar_tool.services.solar_tool.cos')
    def test_calculate_ranges(self, cos_mock):
        # Test case with sample inputs
        longitude = -122.4194
        latitude = 37.7749
        expected_longitude_range = (-122.44, -122.399)
        expected_latitude_range = (37.758, 37.792)

        # Mock the cos function
        cos_mock.return_value = 0.813065222

        # Call the function
        longi_range, lati_range = self.instance.calculate_ranges(longitude, latitude)

        # Verify the output
        assert longi_range == expected_longitude_range
        assert lati_range == expected_latitude_range

        # Verify that cos was called with the correct argument
        cos_mock.assert_called_once_with(radians(latitude))

    def test_calculate_ranges_near_poles(self):
        # Test case with input near the poles
        longitude = -122.4194
        latitude = 89.0
        expected_longitude_range = (-123.386, -121.453)
        expected_latitude_range = (88.983, 89.017)

        # Call the function
        longi_range, lati_range = self.instance.calculate_ranges(longitude, latitude)

        # Verify the output
        assert longi_range == expected_longitude_range
        assert lati_range == expected_latitude_range

    def test_calculate_ranges_near_equator(self):
        # Test case with input near the equator
        longitude = -122.4194
        latitude = 0
        expected_longitude_range = (-122.436, -122.403)
        expected_latitude_range = (-0.017, 0.017)

        # Call the function
        longi_range, lati_range = self.instance.calculate_ranges(longitude, latitude)

        # Verify the output
        assert longi_range == expected_longitude_range
        assert lati_range == expected_latitude_range


    def test_search_saved_files(self, monkeypatch):
        # Define some test data
        longitude_range = (23.558, 23.576)
        latitude_range = (43.868, 43.884)
        expected_result = "23.567_43.876.json"
        
        # Mock the listdir function to return the expected file
        mock_listdir = MagicMock(return_value=[expected_result])
        monkeypatch.setattr("os.listdir", mock_listdir)

        # Call the function
        result = self.instance.search_saved_files(longitude_range, latitude_range)

        # Assert the result
        assert result == expected_result
        
        # Assert that the listdir function was called with the correct argument
        mock_listdir.assert_called_once_with(self.instance.SOLAR_DATA_DIR)
    
    def test_search_saved_files_no_files(self, monkeypatch):
        # Define some test data
        longitude_range = (23.558, 23.576)
        latitude_range = (43.868, 43.884)
        expected_result = None
        
        # Mock the listdir function to return an empty list
        mock_listdir = MagicMock(return_value=[])
        monkeypatch.setattr("os.listdir", mock_listdir)

        # Call the function
        result = self.instance.search_saved_files(longitude_range, latitude_range)

        # Assert the result
        assert result == expected_result
        
        # Assert that the listdir function was called with the correct argument
        mock_listdir.assert_called_once_with(self.instance.SOLAR_DATA_DIR)
        
    def test_search_saved_files_multiple_files(self, monkeypatch):
        # Define some test data
        longitude_range = (23.558, 23.576)
        latitude_range = (43.868, 43.884)
        expected_result = "23.567_43.876.json"
        unexpected_result = "23.569_43.880.json"
        
        # Mock the listdir function to return both files
        mock_listdir = MagicMock(return_value=[expected_result, unexpected_result])
        monkeypatch.setattr("os.listdir", mock_listdir)

        # Call the function
        result = self.instance.search_saved_files(longitude_range, latitude_range)

        # Assert the result
        assert result == expected_result
        
        # Assert that the listdir function was called with the correct argument
        mock_listdir.assert_called_once_with(self.instance.SOLAR_DATA_DIR)
        
    def test_search_saved_files_no_matching_files(self, monkeypatch):
        # Define some test data
        longitude_range = (23.558, 23.576)
        latitude_range = (43.868, 43.884)
        expected_result = None
        unexpected_result = "23.589_43.890.json"
        
        # Mock the listdir function to return the unexpected file
        mock_listdir = MagicMock(return_value=[unexpected_result])
        monkeypatch.setattr("os.listdir", mock_listdir)

        # Call the function
        result = self.instance.search_saved_files(longitude_range, latitude_range)

        # Assert the result
        assert result == expected_result
        
        # Assert that the listdir function was called with the correct argument
        mock_listdir.assert_called_once_with(self.instance.SOLAR_DATA_DIR)


    def test_save_json_into_file(self, monkeypatch):
        # Define some test data
        output_data = {"longitude": 23.567, "latitude": 43.876, "solar_data": {"irradiance": 1015.0, "temperature": 23.5}}
        output_filename = "23.567_43.876.json"
        
        # Mock the open function to return a file object
        mock_file = mock_open()
        monkeypatch.setattr("builtins.open", mock_file)

        # Call the function
        self.instance.save_json_into_file(output_data, output_filename)

        # Assert that the file was created
        mock_file.assert_called_once_with(self.instance.SOLAR_DATA_DIR + "/" + output_filename, 'w')
        
        # Assert that the dump function was called with the correct arguments
        mock_file().write.call_count == 21

    def test_save_json_into_file_os_error(self, monkeypatch):
        # Define some test data
        output_data = {"longitude": 23.567, "latitude": 43.876, "solar_data": {"irradiance": 1015.0, "temperature": 23.5}}
        output_filename = "23.567_43.876.json"
        
        # Mock the open function to raise an OSError
        mock_file = mock_open()
        mock_file.side_effect = OSError
        monkeypatch.setattr("builtins.open", mock_file)

        # Call the function
        with pytest.raises(OSError):
            self.instance.save_json_into_file(output_data, output_filename)

        # Assert that the file was not created
        mock_file.assert_called_once_with(self.instance.SOLAR_DATA_DIR + "/" + output_filename, 'w')

    @patch.object(SolarTool, 'calculate_ranges')
    @patch.object(SolarTool, 'search_saved_files')
    def test_get_solar_radiation_file_found(self, mock_search_saved_files, mock_calculate_ranges):
        # Create a mock SolarDataProcessor instance
        mock_calculate_ranges.return_value = ((-100, -99), (40, 41))
        mock_search_saved_files.return_value = "test_file.json"

        # Create a mock file object to simulate reading data from a file
        file_mock = MagicMock()
        file_mock.__enter__.return_value = file_mock
        file_mock.read.return_value = json.dumps({"test_data": "test_value"})

        # Patch the built-in 'open' method to return the mock file object
        with patch("builtins.open", return_value=file_mock):
            # Call the method being tested
            result = self.instance.get_solar_radiation(-99.5, 40.5)

        # Ensure that the file was opened and read
        assert file_mock.__enter__.called
        assert file_mock.read.called

        # Ensure that the correct data was returned
        assert result == {"test_data": "test_value"}

    @patch.object(SolarTool, 'get_3day_moving_avg')
    @patch.object(SolarTool, 'get_solar_radiation_from_api')
    @patch.object(SolarTool, 'calculate_ranges')
    @patch.object(SolarTool, 'search_saved_files')
    def test_get_solar_radiation_file_not_found(self, mock_search_saved_files, mock_calculate_ranges, mock_get_solar_radiation_from_api, mock_get_3day_moving_avg):
        # Create a mock SolarDataProcessor instance
        mock_calculate_ranges.return_value = ((-100, -99), (40, 41))
        mock_search_saved_files.return_value = None
        mock_get_solar_radiation_from_api.return_value = {
            "properties": {
                "parameter": {
                    "ALLSKY_SFC_SW_DWN": {
                        "20210101": 1.0,
                        "20210102": -999,
                        "20210103": 3.0
                    }
                }
            }
        }
        mock_get_3day_moving_avg.return_value = {
            "20210103": 2.0
        }

        # Create a mock file object to simulate writing data to a file
        file_mock = MagicMock()
        file_mock.__enter__.return_value = file_mock
        file_mock.write.return_value = None

        # Patch the built-in 'open' method to return the mock file object
        with patch("builtins.open", return_value=file_mock):
            # Call the method being tested
            result = self.instance.get_solar_radiation(-99.5, 40.5)

        # Ensure that the correct data was returned
        assert result == {
            "properties": {
                "parameter": {
                    "ALLSKY_SFC_SW_DWN": {
                        "20210101": 1.0,
                        "20210103": 3.0
                    },
                    "ALLSKY_SFC_SW_DWN_MOVING_AVG": {
                        "20210103": 2.0
                    }
                }
            }
        }
    
    @patch('requests.get')
    def test_get_solar_radiation_from_api_success(self, mock_get):
        # arrange
        expected_longitude = 100
        expected_latitude = 50
        mock_response = Mock()
        mock_response.status_code = 200
        mock_response.json.return_value = {'properties': {'parameter': {'ALLSKY_SFC_SW_DWN': {'20230101': 10}}}}
        mock_get.return_value = mock_response

        # act
        result = self.instance.get_solar_radiation_from_api(expected_longitude, expected_latitude)

        # assert
        mock_get.assert_called_once_with(f"https://power.larc.nasa.gov/api/temporal/daily/point?parameters=ALLSKY_SFC_SW_DWN&community=RE&longitude={expected_longitude}&latitude={expected_latitude}&start=19840101&end=20221231&format=JSON")
        assert result == {'properties': {'parameter': {'ALLSKY_SFC_SW_DWN': {'20230101': 10}}}}

    @patch('requests.get')
    def test_get_solar_radiation_from_api_failure(self, mock_get):
        # arrange
        expected_longitude = 100
        expected_latitude = 50
        mock_response = Mock()
        mock_response.status_code = 500
        mock_get.return_value = mock_response

        # act & assert
        with pytest.raises(Exception) as excinfo:
            self.instance.get_solar_radiation_from_api(expected_longitude, expected_latitude)

        assert str(excinfo.value).startswith("Error: 500")


    def test_empty_data(self):
        data = {"properties": {"parameter": {"ALLSKY_SFC_SW_DWN": {}}}}
        result = self.instance.get_3day_moving_avg(data)
        assert result == {}

    def test_one_day_data(self):
        data = {"properties": {"parameter": {"ALLSKY_SFC_SW_DWN": {"20220101": 10}}}}
        result = self.instance.get_3day_moving_avg(data)
        assert result == {"20220101": 10}

    def test_two_day_data(self):
        data = {"properties": {"parameter": {"ALLSKY_SFC_SW_DWN": {"20220101": 10, "20220102": 20}}}}
        result = self.instance.get_3day_moving_avg(data)
        assert result == {"20220101": 10, "20220102": 15}

    def test_three_day_data(self):
        data = {"properties": {"parameter": {"ALLSKY_SFC_SW_DWN": {"20220101": 10, "20220102": 20, "20220103": 30}}}}
        result = self.instance.get_3day_moving_avg(data)
        assert result == {"20220101": 10, "20220102": 15.0, "20220103": 20.0}

    def test_four_day_data(self):
        data = {"properties": {"parameter": {"ALLSKY_SFC_SW_DWN": {"20220101": 10, "20220102": 20, "20220103": 30, "20220104": 40}}}}
        result = self.instance.get_3day_moving_avg(data)
        assert result == {"20220101": 10, "20220102": 15.0, "20220103": 20.0, "20220104": 30.0}

    def test_invalid_data(self):
        data = {"properties": {"parameter": {"ALLSKY_SFC_SW_DWN": {"20220101": "invalid", "20220102": 20, "20220103": 30, "20220104": 40}}}}
        with pytest.raises(TypeError):
            self.instance.get_3day_moving_avg(data)
    
    @pytest.fixture
    def mock_get_solar_radiation(self):
        """Mock get_solar_radiation function"""
        return MagicMock(return_value={
            'properties': {
                'parameter': {
                    'ALLSKY_SFC_SW_DWN_MOVING_AVG': {
                        '20230101': 3.0,
                        '20230102': 4.0,
                        '20230103': 5.0
                    },
                    'ALLSKY_SFC_SW_DWN': {
                        '20230101': 6.0,
                        '20230102': 7.0,
                        '20230103': 8.0
                    }
                }
            }
        })

    def test_predict_solar_tools_requirements(self, mock_get_solar_radiation):
        """Test predict_solar_tools_requirements function"""
        site =  { "longitude": 20, "latitude": 30, "power_budget": 1000, "reliability_target": 90 }

        # Mock the get_solar_radiation function
        self.instance.get_solar_radiation = mock_get_solar_radiation

        # Test the function with the given inputs
        with patch('app.solar_tool.services.solar_tool.subprocess.check_call', return_value=''):
            result = self.instance.predict_solar_tools_requirements(site)

        # Check the expected output
        assert result.number_of_solar_modules == 23
        assert result.solar_pv_to_install_watts == 9200
        assert result.number_of_batteries == 18
        assert result.batteries_capacity_to_install_kWh == 43.2
        assert result.max_output_angle == "30 degrees south"

    def test_predict_solar_tools_requirements_exception(self, mock_get_solar_radiation):
        """Test predict_solar_tools_requirements function with an exception"""
        site = { "longitude": 20, "latitude": 30, "power_budget": 1000, "reliability_target": 90 }

        # Mock the get_solar_radiation function
        self.instance.get_solar_radiation = mock_get_solar_radiation

        # Raise an exception to test the error handling
        mock_get_solar_radiation.side_effect = Exception("Failed to get solar radiation data")

        # Test the function with the given inputs and expect an exception
        with pytest.raises(Exception) as ex:
            with patch('app.solar_tool.services.solar_tool.subprocess.check_call', return_value=''):
                self.instance.predict_solar_tools_requirements(site)

        # Check the expected exception message
        assert str(ex.value) == "Failed to get solar radiation data"

    def test_predict_solar_tools_requirements_zero_power_budget(self, mock_get_solar_radiation):
        """Test predict_solar_tools_requirements function with a zero power budget"""
        site = { "longitude": 20, "latitude": 30, "power_budget": 0, "reliability_target": 90 }

        # Mock the get_solar_radiation function
        self.instance.get_solar_radiation = mock_get_solar_radiation

        # Test the function with zero power budget
        with patch('app.solar_tool.services.solar_tool.subprocess.check_call', return_value=''):
            result = self.instance.predict_solar_tools_requirements(site)

        # Check the expected output
        assert result.number_of_solar_modules == 0
        assert result.solar_pv_to_install_watts == 0
        assert result.number_of_batteries == 0
        assert result.batteries_capacity_to_install_kWh == 0.0