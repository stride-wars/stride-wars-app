// import React from 'react'
// import { Screen } from '../App'
// import { MapIcon, PlayIcon, SettingsIcon, UserIcon } from 'lucide-react'
// interface HeaderProps {
//   currentScreen: Screen
//   onNavigate: (screen: Screen) => void
//   username: string
// }
// export const Header: React.FC<HeaderProps> = ({
//   currentScreen,
//   onNavigate,
//   username,
// }) => {
//   return (
//     <header className="w-full bg-gray-800 p-4">
//       <div className="flex items-center justify-between">
//         <div className="flex items-center">
//           <img
//             src="https://uploadthingy.s3.us-west-1.amazonaws.com/qqRNpchVnNPFtnb8JGxVvX/20250416_1813_Stride_Wars_Logo_simple_compose_01jrznwr8ke9bv20fhxtvk166n%281%29.png"
//             alt="Stride Wars Logo"
//             className="h-10 mr-2"
//           />
//           <h1 className="text-xl font-bold text-yellow-500">Stride Wars</h1>
//         </div>
//         <div className="flex items-center">
//           <span className="mr-2 text-sm hidden md:inline">{username}</span>
//           <UserIcon className="h-5 w-5 text-yellow-500" />
//         </div>
//       </div>
//       <nav className="mt-4">
//         <ul className="flex justify-around">
//           <li>
//             <button
//               onClick={() => onNavigate(Screen.MAP)}
//               className={`flex flex-col items-center p-2 ${currentScreen === Screen.MAP ? 'text-yellow-500' : 'text-gray-400'}`}
//             >
//               <MapIcon className="h-6 w-6" />
//               <span className="text-xs mt-1">Map</span>
//             </button>
//           </li>
//           <li>
//             <button
//               onClick={() => onNavigate(Screen.ACTIVITY)}
//               className={`flex flex-col items-center p-2 ${currentScreen === Screen.ACTIVITY ? 'text-yellow-500' : 'text-gray-400'}`}
//             >
//               <PlayIcon className="h-6 w-6" />
//               <span className="text-xs mt-1">Activity</span>
//             </button>
//           </li>
//           <li>
//             <button
//               onClick={() => onNavigate(Screen.SETTINGS)}
//               className={`flex flex-col items-center p-2 ${currentScreen === Screen.SETTINGS ? 'text-yellow-500' : 'text-gray-400'}`}
//             >
//               <SettingsIcon className="h-6 w-6" />
//               <span className="text-xs mt-1">Settings</span>
//             </button>
//           </li>
//         </ul>
//       </nav>
//     </header>
//   )
// }
